package main

import (
	"bufio"
	"log"
	"os"
	"strconv"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	channel = 0
	name    = "drive-test"
)

var (
	frequency = 50
	minPulse  = 750
	maxPulse  = 2250
)

// ESC holds all data
type ESC struct {
	device *i2c.PCA9685Driver

	_throttle float32 // [-1.0,+1.0]
	_fraction float32 // [0,1]

	pulse    uint16
	minPulse uint16
	maxPulse uint16
}

// NewSpeedController creates a new speed controller
func NewSpeedController(d *i2c.PCA9685Driver, min, max int) (*ESC, error) {
	esc := ESC{
		device:    d,
		_throttle: 0.0,
		_fraction: 0.5,
		minPulse:  uint16(min),
		maxPulse:  uint16(max),
	}

	return &esc, nil
}

func (esc *ESC) toString() {
	log.Printf("%v", esc)
}

func (esc *ESC) throttle(value float32) {
	if value < -1.0 {
		esc._throttle = -1.0
	} else if value > 1.0 {
		esc._throttle = 1.0
	} else {
		esc._throttle = value
	}

	esc._fraction = (esc._throttle + 1.0) / 2.0
	esc.pulse = esc.minPulse + uint16((float32(esc.maxPulse-esc.minPulse) * esc._fraction))

	err := esc.device.SetPWM(channel, uint16(0), esc.pulse)
	if err != nil {
		log.Printf(err.Error())
	}
}

func (esc *ESC) setThrottle(value uint16) {
	if value > 4096 {
		value = 4096
	}

	esc._fraction = float32(value) / 4096.0
	esc.pulse = value
	err := esc.device.SetPWM(channel, uint16(0), value)
	if err != nil {
		log.Printf(err.Error())
	}
}

func drive(esc *ESC, start, stop int) error {
	log.Printf("Drive Test Run Loop...\n")

	reader := bufio.NewReader(os.Stdin)

	increment := 5

	for pulse := start + increment; pulse > start; {

		char, _, _ := reader.ReadRune()
		switch char {
		case 'w':
			pulse = pulse + increment
			esc.device.SetPWM(channel, uint16(start), uint16(pulse))
			log.Printf("%d %d\n", start, pulse)

			break
		case 's':
			pulse = pulse - increment
			esc.device.SetPWM(channel, uint16(start), uint16(pulse))
			log.Printf("%d %d\n", start, pulse)

			break
		}
	}

	esc.setThrottle(0)
	return nil
}

/*

320/580 -> beep
320/590 -> first noise
320/605 -> slow movement
320/710 -> max

*/

func main() {

	// read some parameters from the command line
	if len(os.Args) > 0 {

		if len(os.Args) == 4 {
			frequency, _ = strconv.Atoi(os.Args[1])
			// servo pulse range
			minPulse, _ = strconv.Atoi(os.Args[2])
			maxPulse, _ = strconv.Atoi(os.Args[3])
		}
	}

	r := raspi.NewAdaptor()
	pca9685 := i2c.NewPCA9685Driver(r)

	work := func() {
		err := pca9685.SetPWMFreq(float32(frequency))
		if err != nil {
			log.Printf(err.Error())
			return
		}

		esc, err := NewSpeedController(pca9685, minPulse, maxPulse)
		if err != nil {
			log.Printf(err.Error())
			return
		}

		drive(esc, minPulse, maxPulse)
	}

	robot := gobot.NewRobot(name,
		[]gobot.Connection{r},
		[]gobot.Device{pca9685},
		work,
	)

	robot.Start()
}
