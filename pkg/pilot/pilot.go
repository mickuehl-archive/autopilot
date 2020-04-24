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
)

type (
	// Config holds basic pre-sets
	Config struct {
		frequency int // PCA9685 clock speed
		// actuators
		servoChan int // channel of the servo
		escChan   int // channel of the speed controller
		led1Chan  int // channels of the two back/break LEDs
		led2Chan  int

		// actuator pre-sets

		// min pulse length out of 4096
		servoMinPulse int
		// max pulse length out of 4096
		servoMaxPulse int
		// the max the servo can rotate (in deg)
		servoMaxRange int
		// the max the servo is allowed to rotate to each side (deg)
		servoRangeLimit int

		// min pulse length out of 4096, rel to escZeroPulse
		escMinPulse int
		// max pulse length out of 4096, rel to escZeroPulse
		escMaxPulse int
		// zero pulse length out of 4096
		escZeroPulse int
		// base pulse to work from
		escBasePulse int
		// is the pulse to reset/initialize the ESC
		escInitPulse int
		// esc can reverse if true
		escCanReverse bool
	}

	// Pilot is a collection of assets that make up a vehicle
	Pilot struct {
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
	}
)

// NewPilot creates a new instance of the OBU controlling the vehicle
func NewPilot() (*Pilot, error) {

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
	p := &Pilot{
		cfg:       cfg,
		adaptor:   r,
		actuators: pca9685,
		steering:  0,
		throttle:  0.0,
	}

	return p, nil
}

// Start initializes the pilot and all its components
func (p *Pilot) Start() error {
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
	time.Sleep(INIT_DELAY * time.Millisecond)
	// all good
	return nil
}

// Direction sets the steering angle
func (p *Pilot) Direction(value int) {
	if value < p.cfg.servoRangeLimit*-1 {
		p.steering = p.cfg.servoRangeLimit * -1
	} else if value > p.cfg.servoRangeLimit {
		p.steering = p.cfg.servoRangeLimit
	} else {
		p.steering = value
	}

	// value == 0 -> servoMaxRange / 2
	direction := (p.cfg.servoMaxRange / 2) + p.steering + p.steeringTrim
	pulse := uint16(degree2pulse(direction, p.cfg.servoMinPulse, p.cfg.servoMaxPulse, p.cfg.servoMaxRange))

	// set the servo pulse
	err := p.actuators.SetPWM(p.cfg.servoChan, uint16(0), pulse)
	if err != nil {
		log.Printf(err.Error())
	}

}

// Throttle sets the motor speed
func (p *Pilot) Throttle(value float32) {

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
	log.Printf("Throttle: %f => (%d,%d)", p.throttle, p.cfg.escBasePulse, pulseOff) // FIXME remove this
	err := p.actuators.SetPWM(p.cfg.escChan, uint16(p.cfg.escBasePulse), uint16(pulseOff))
	if err != nil {
		log.Printf(err.Error())
	}

}

// Shutdown stops all components
func (p *Pilot) Shutdown() error {
	log.Println("Shutdown")

	// set the actuators to a neutral position
	p.Direction(0)
	p.Throttle(0.0)

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

func degree2pulse(deg, min, max, r int) int {
	return min + ((max-min)/r)*deg
}

func throttle2pulse(throttle float32, zero, min, max int, reverse bool) int {
	// map the throttle [-1.0,1.0] to [0,1.0]
	return zero + int(throttle*float32(max-min))
}

// map a value from range r1 onto range r2
func mapRange(x, r1min, r1max, r2min, r2max int) int {
	ratio := float64(r2max-r2min) / float64(r1max-r1min)
	return int(float64(x) * ratio)
}
