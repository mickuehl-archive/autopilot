package parts

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/majordomusio/commons/pkg/util"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"

	"shadow-racer/autopilot/v1/pkg/eventbus"
	"shadow-racer/autopilot/v1/pkg/obu"
)

const (
	frequency      = 50
	blinkFrequency = "220ms"
	// channel assignments
	throttleChan = 0
	steeringChan = 3
	led1Chan     = 8
	led2Chan     = 11

	stateDriving = "DRIVING"
	stateStopped = "STOPPED"
)

type (
	// RaspiOnboardUnit holds data needed to simulate an OBU
	RaspiOnboardUnit struct {
		mutex *sync.Mutex
		// hardware
		adaptor   *raspi.Adaptor     // the Raspi
		actuators *i2c.PCA9685Driver // To controll actuators like servo, ESC, LEDs
		// actuators
		servo *StandardServo
		esc   *StandardSpeedController
		// some other state
		bkLights int // 0 .. 4000. 0 = off, 4000 = max

		// high-level vehicle state
		vehicle *obu.Vehicle
	}
)

// NewRaspiOnboardUnit creates a new instance of a real OBU including actuators etc.
func NewRaspiOnboardUnit() (*RaspiOnboardUnit, error) {

	// a Raspberry Pi as the platform
	r := raspi.NewAdaptor()
	// a board with a PCA9685 to control servos
	pca9685 := i2c.NewPCA9685Driver(r)
	if pca9685 == nil {
		log.Fatalf("Could not initialize the PCA9685 driver")
		return nil, errors.New("Could not initialize the PCA9685 driver")
	}
	pca9685.SetName("pca9685")

	// pull it all together
	obu := &RaspiOnboardUnit{
		mutex:     &sync.Mutex{},
		adaptor:   r,
		actuators: pca9685,
		servo:     NewBMS390DMH(steeringChan),
		esc:       NewWP40(throttleChan),
		vehicle: &obu.Vehicle{
			Mode:        stateStopped,
			Throttle:    0.0,
			Steering:    0.0,
			Heading:     0.0,
			Recording:   false,
			RecordingTS: 0,
			TS:          util.TimestampNano(),
		},
	}
	// set a speed limit for now
	obu.esc.Limit = 30

	return obu, nil
}

// Initialize prepares the device
func (o *RaspiOnboardUnit) Initialize() error {

	err := o.actuators.Start()
	if err != nil {
		return err
	}

	// started all components, wait a bit before further configuration happens ...
	time.Sleep(500 * time.Millisecond)
	o.actuators.SetPWMFreq(float32(frequency))

	// calibrate & reset the esc
	if o.esc.Cfg.InitPulse > 0 {
		o.SetChannelPulse(o.esc.Cfg.N, o.esc.Cfg.BasePulse, o.esc.Cfg.InitPulse)
		time.Sleep(500 * time.Millisecond)
		o.SetChannelPulse(o.esc.Cfg.N, o.esc.Cfg.BasePulse, o.esc.Cfg.ZeroPulse)
		time.Sleep(500 * time.Millisecond)
	}

	// start all event handlers
	go o.RCStateHandler()

	// all good
	logger.Info("OBU is ready")
	return nil
}

// Reset re-initializes the device
func (o *RaspiOnboardUnit) Reset() error {
	logger.Debug("reset the raspi")
	return nil
}

// Shutdown releases all resources
func (o *RaspiOnboardUnit) Shutdown() error {
	// stop the hardware
	o.actuators.Halt()
	return o.adaptor.Finalize()
}

// Event handlers etc

// RCStateHandler listens on events for remote state changes
func (o *RaspiOnboardUnit) RCStateHandler() {
	ch := eventbus.InstanceOf().Subscribe("rc/state")
	for {
		evt := <-ch
		state := evt.Data.(RemoteState)

		o.mutex.Lock()

		if state.Mode != o.vehicle.Mode {
			if state.Mode == stateDriving {
				// assumes o.vehicle.Mode == STOPPED
				o.TailLights(4000, true)
			} else if state.Mode == stateStopped {
				o.TailLightsOff()
				o.vehicle.Throttle = 0.0
				o.vehicle.Steering = 0.0
			} else {
				// FIXME should not happen
			}
			o.vehicle.Mode = state.Mode
			o.vehicle.Recording = state.Recording
		} else {
			o.vehicle.Steering = 100.0 * ((float32(o.servo.MaxRange) / 90.0) * state.Steering)
			o.vehicle.Throttle = 100.0 * state.Throttle
		}

		if state.Recording != o.vehicle.Recording {
			baseURL := "http://localhost:3001" // FIXME configuration

			if state.Recording == true {
				o.vehicle.RecordingTS = util.Timestamp()
				o.vehicle.Recording = true

				resp, err := http.Get(fmt.Sprintf("%s/start?ts=%d", baseURL, o.vehicle.RecordingTS))

				if err != nil {
					logger.Error("Error toggling recording", "err", err.Error())
				} else {
					logger.Info("Started recording", "ts", o.vehicle.RecordingTS)
				}
				defer resp.Body.Close()
			} else {
				o.vehicle.Recording = false
				resp, err := http.Get(baseURL + "/stop")

				if err != nil {
					logger.Error("Error toggling recording", "err", err.Error())
				} else {
					logger.Info("Stopped recording")
				}
				defer resp.Body.Close()
			}
		}

		o.vehicle.TS = util.TimestampNano()

		// publish the new state
		eventbus.InstanceOf().Publish("state/vehicle", o.vehicle.Clone())

		// set the actuators
		o.Direction(int(o.vehicle.Steering))
		o.Throttle(int(o.vehicle.Throttle))

		o.mutex.Unlock()
	}
}

// OBU specific functions

// SetChannelPulse sets the pulse values for a channel
func (o *RaspiOnboardUnit) SetChannelPulse(ch, pulseOn, pulseOff int) error {
	//logger.Debug("Channel on/off pulse", "channel", ch, "on", pulseOn, "off", pulseOff)

	if ch < 0 || ch > 15 {
		return errors.New("Invalid channel")
	}
	if pulseOn < 0 || pulseOn > 4096 {
		return errors.New("Invalid pulse 'on' value")
	}
	if pulseOff < 0 || pulseOff > 4096 {
		return errors.New("Invalid pulse 'off' value")
	}
	return o.actuators.SetPWM(ch, uint16(pulseOn), uint16(pulseOff))
}

// Direction sets the steering direction (in deg)
func (o *RaspiOnboardUnit) Direction(value int) {
	// expect servo to calculate the pulse values
	on, off := o.servo.SetAngle(value)
	// set the values on the channel
	o.SetChannelPulse(o.servo.Cfg.N, on, off)
}

// Throttle sets the speed (0..100)
func (o *RaspiOnboardUnit) Throttle(value int) {
	// expect ESC to calculate the pulse values
	on, off := o.esc.SetThrottle(value)
	// set the values on the channel
	o.SetChannelPulse(o.esc.Cfg.N, on, off)
}

// TailLights sets the taillights/brake lights (value = 0 off, value = 4000 max)
func (o *RaspiOnboardUnit) TailLights(value int, blink bool) {
	if value < 0 {
		o.bkLights = 0
	} else if value > 4000 {
		o.bkLights = 4000
	} else {
		o.bkLights = value
	}
	if blink {
		// blink in a go routine
		go func() {
			pause, err := time.ParseDuration(blinkFrequency)
			if err != nil {
				return
			}
			for true {
				o.SetChannelPulse(led1Chan, 0, o.bkLights)
				o.SetChannelPulse(led2Chan, 0, o.bkLights)
				time.Sleep(pause)
				o.SetChannelPulse(led1Chan, 0, 0)
				o.SetChannelPulse(led2Chan, 0, 0)
				time.Sleep(pause)

				if o.bkLights == 0 {
					break
				}
			}
		}()
	} else {
		o.SetChannelPulse(led1Chan, 0, o.bkLights)
		o.SetChannelPulse(led2Chan, 0, o.bkLights)
	}
}

// TailLightsOff turns the taillights/brake lights off
func (o *RaspiOnboardUnit) TailLightsOff() {
	o.bkLights = 0
	o.SetChannelPulse(led1Chan, 0, 0)
	o.SetChannelPulse(led2Chan, 0, 0)
}
