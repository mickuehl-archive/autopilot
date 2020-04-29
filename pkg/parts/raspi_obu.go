package parts

import (
	"errors"
	"log"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	frequency      = 50
	blinkFrequency = "220ms"
	// channel assignments
	led1Chan = 8
	led2Chan = 11
)

type (
	// RaspiOnboardUnit holds data needed to simulate an OBU
	RaspiOnboardUnit struct {
		// hardware
		adaptor   *raspi.Adaptor     // the Raspi
		actuators *i2c.PCA9685Driver // To controll actuators like servo, ESC, LEDs
		// state
		bkLights int // 0 .. 4000. 0 = off, 4000 = max
	}
)

// NewRaspiOnboardUnit creates a new instance of a virtual OBU, for e.g. unit testing of the framework
func NewRaspiOnboardUnit() (*RaspiOnboardUnit, error) {

	// a Raspberry Pi as the platform
	r := raspi.NewAdaptor()
	// a board with a PCA9685 to control servos etc
	pca9685 := i2c.NewPCA9685Driver(r)
	if pca9685 == nil {
		log.Fatalf("Could not initialize the PCA9685 driver")
		return nil, errors.New("Could not initialize the PCA9685 driver")
	}
	pca9685.SetName("pca9685")

	obu := &RaspiOnboardUnit{
		adaptor:   r,
		actuators: pca9685,
	}
	return obu, nil
}

// Initialize prepares the device
func (obu *RaspiOnboardUnit) Initialize() error {
	logger.Debug("initialize the raspi")

	err := obu.actuators.Start()
	if err != nil {
		return err
	}

	// started all components, wait a bit before further configuration happens ...
	time.Sleep(500 * time.Millisecond)
	obu.actuators.SetPWMFreq(float32(frequency))

	return nil
}

// Reset re-initializes the device
func (obu *RaspiOnboardUnit) Reset() error {
	logger.Debug("reset the raspi")
	return nil
}

// Shutdown releases all resources
func (obu *RaspiOnboardUnit) Shutdown() error {
	logger.Debug("shutdown the raspi")

	// stop the hardware
	obu.actuators.Halt()
	return obu.adaptor.Finalize()
}

// OBU specific functions

// TailLights sets the taillights/brake lights (value = 0 off, value = 4000 max)
func (obu *RaspiOnboardUnit) TailLights(value int, blink bool) {
	if value < 0 {
		obu.bkLights = 0
	} else if value > 4000 {
		obu.bkLights = 4000
	} else {
		obu.bkLights = value
	}
	if blink {
		// blink in a go routine
		go func() {
			pause, err := time.ParseDuration(blinkFrequency)
			if err != nil {
				log.Printf(err.Error())
				return
			}
			for true {
				obu.actuators.SetPWM(led1Chan, 0, uint16(obu.bkLights))
				obu.actuators.SetPWM(led2Chan, 0, uint16(obu.bkLights))
				time.Sleep(pause)
				obu.actuators.SetPWM(led1Chan, 0, 0)
				obu.actuators.SetPWM(led2Chan, 0, 0)
				time.Sleep(pause)
				if obu.bkLights == 0 {
					break
				}
			}
		}()
	} else {
		obu.actuators.SetPWM(led1Chan, 0, uint16(obu.bkLights))
		obu.actuators.SetPWM(led2Chan, 0, uint16(obu.bkLights))
	}
}

// TailLightsOff turns the taillights/brake lights off
func (obu *RaspiOnboardUnit) TailLightsOff() {
	obu.bkLights = 0
	obu.actuators.SetPWM(led1Chan, 0, 0)
	obu.actuators.SetPWM(led2Chan, 0, 0)
}
