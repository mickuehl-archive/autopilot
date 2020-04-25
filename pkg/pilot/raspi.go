package pilot

import (
	"errors"
	"log"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	// INIT_DELAY is the time to wait until settings apply
	INIT_DELAY = 300
	BLINK_FREQ = "220ms"
)

type (
	// Pilot is a collection of assets that make up a vehicle
	RaspiPilot struct {
		// configuration
		cfg *Config
		// hardware
		adaptor   *raspi.Adaptor     // the Raspi
		actuators *i2c.PCA9685Driver // servo, ESC, LEDs
		// state information
		steering     int     // direction of the steering (in deg), 0 is center. Range is [-servoRangeLimit,servoRangeLimit]
		steeringTrim int     // steering correction factor
		throttle     float32 // motor speed, -1.0 (full reverse) to 1.0 (full speed)
		throttleTrim float32 // motor speed trim
		bkLights     int     // 0 .. 4000
	}
)

// NewRaspiPilot creates a new instance of the Raspi OBU controlling the vehicle
func NewRaspiPilot() (*RaspiPilot, error) {

	// get/load the basic configuration
	cfg := newConfig()
	if cfg == nil {
		log.Fatalf("Could not initialize the global configuration")
		return nil, errors.New("Could not initialize the global configuration")
	}

	// a Raspberry Pi as the platform
	r := raspi.NewAdaptor()
	// a board with a PCA9685 to control servos
	pca9685 := i2c.NewPCA9685Driver(r)
	if pca9685 == nil {
		log.Fatalf("Could not initialize the PCA9685 driver")
		return nil, errors.New("Could not initialize the PCA9685 driver")
	}
	pca9685.SetName("pca9685")

	// the main data structure
	p := &RaspiPilot{
		cfg:       cfg,
		adaptor:   r,
		actuators: pca9685,
		steering:  0,
		throttle:  0.0,
		bkLights:  0,
	}

	return p, nil
}

// Start initializes the pilot and all its components
func (p *RaspiPilot) Start() error {
	log.Println("Starting")

	err := p.actuators.Start()
	if err != nil {
		return err
	}

	// started all components, wait a bit before further configuration happens ...
	time.Sleep(500 * time.Millisecond)
	p.actuators.SetPWMFreq(float32(p.cfg.frequency))

	// wait for the hardware to be ready then calibrate the vehicle
	time.Sleep(1000 * time.Millisecond)
	// steering to zero
	p.Direction(0)
	// calibrate & reset the esc
	p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(p.cfg.escInitPulse))
	time.Sleep(INIT_DELAY * time.Millisecond)
	p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(p.cfg.escZeroPulse))
	time.Sleep(INIT_DELAY * time.Millisecond)
	p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(p.cfg.escZeroPulse-20))

	// final delay
	time.Sleep(1 * time.Second)
	// all good
	return nil
}

// Direction sets the steering angle
func (p *RaspiPilot) Direction(value int) {
	if value < p.cfg.servoRangeLimit*-1 {
		p.steering = p.cfg.servoRangeLimit * -1
	} else if value > p.cfg.servoRangeLimit {
		p.steering = p.cfg.servoRangeLimit
	} else {
		p.steering = value
	}

	// value == 0 -> servoMaxRange / 2
	direction := (p.cfg.servoMaxRange / 2) + p.steering + p.steeringTrim
	pulse := p.cfg.servoMinPulse + ((p.cfg.servoMaxPulse-p.cfg.servoMinPulse)/p.cfg.servoMaxRange)*direction

	// set the servo pulse
	//log.Printf("Direction: %d => (%d,%d)", p.steering, 0, pulse) // FIXME remove this
	err := p.actuators.SetPWM(p.cfg.servoChan, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}

}

// Throttle sets the motor speed
func (p *RaspiPilot) Throttle(value float32) {

	changeOfDirection := false

	if p.cfg.escCanReverse {
		if (p.throttle > 0.0 && value <= 0.0) || (p.throttle < 0.0 && value >= 0.0) {
			changeOfDirection = true
		}

		// can reverse
		if value < -1.0 {
			p.throttle = -1.0
		} else if value > 1.0 {
			p.throttle = 1.0
		} else {
			p.throttle = value
		}
	} else {
		// no reverse
		if value < 0.0 {
			p.throttle = 0.0
		} else if value > 1.0 {
			p.throttle = 1.0
		} else {
			p.throttle = value
		}
	}

	pulseOff := 0
	if p.throttle == 0.0 {
		pulseOff = p.cfg.escZeroPulse
	} else if p.throttle > 0.0 {
		pulseOff = p.cfg.escZeroPulse + int(float32(p.cfg.escMaxPulse)*p.throttle)
	} else {
		pulseOff = p.cfg.escZeroPulse + int(float32(p.cfg.escMinPulse)*p.throttle) // throttle < 0 -> PLUS zeroPulse
	}

	if changeOfDirection {
		p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(p.cfg.escZeroPulse))
		time.Sleep(INIT_DELAY * time.Millisecond)
	}
	// set the servo pulse
	// log.Printf("Throttle: %f => (%d,%d)", p.throttle, p.cfg.escBasePulse, pulseOff) // FIXME remove this
	err := p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(pulseOff))
	if err != nil {
		log.Printf(err.Error())
	}

}

// BackLights turns the back lights on/off, w/o blinking
func (p *RaspiPilot) BackLights(value int, blink bool) {
	if value < 0 {
		p.bkLights = 0
	} else if value > 4000 {
		p.bkLights = 4000
	} else {
		p.bkLights = value
	}
	if blink {
		// blink in a go routine
		go func() {
			pause, err := time.ParseDuration(BLINK_FREQ)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			for true {
				p.actuators.SetPWM(p.cfg.led1Chan, 0, uint16(p.bkLights))
				p.actuators.SetPWM(p.cfg.led2Chan, 0, uint16(p.bkLights))
				time.Sleep(pause)
				p.actuators.SetPWM(p.cfg.led1Chan, 0, 0)
				p.actuators.SetPWM(p.cfg.led2Chan, 0, 0)
				time.Sleep(pause)
				if p.bkLights == 0 {
					break
				}
			}
		}()
	} else {
		p.actuators.SetPWM(p.cfg.led1Chan, 0, uint16(p.bkLights))
		p.actuators.SetPWM(p.cfg.led2Chan, 0, uint16(p.bkLights))
	}

}

// Shutdown stops all components
func (p *RaspiPilot) Shutdown() error {
	log.Println("Shutdown")

	// set the actuators to a neutral position
	p.Direction(0)
	p.Throttle(0.0)
	p.BackLights(0, false)

	// stop the hardware drivers
	p.actuators.Halt()
	return p.adaptor.Finalize()
}

func newConfig() *Config {
	c := Config{
		frequency: 50,
		// channels
		servoChan: 3,
		escChan:   0,
		led1Chan:  8,
		led2Chan:  11,
		// servo
		servoMinPulse:   150,
		servoMaxPulse:   700,
		servoMaxRange:   180,
		servoRangeLimit: 30,
		// ESC
		escMinPulse:   100,
		escMaxPulse:   100,
		escBasePulse:  1000,
		escZeroPulse:  1300,
		escInitPulse:  2000,
		escCanReverse: false, // it can, but does not work reliably though
	}

	return &c
}
