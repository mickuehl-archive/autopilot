package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	frequency      = 50
	channel   byte = 3
)

var (
	// Min pulse length out of 4096
	servoMin = 150
	// Max pulse length out of 4096
	servoMax = 700
	// Limiting the max this servo can rotate (in deg)
	maxDegree = 180
	// Number of degrees to increase per call
	increase = 15
)

func degree2pulse(deg int) int32 {
	pulse := servoMin
	pulse += ((servoMax - servoMin) / maxDegree) * deg
	return int32(pulse)
}

func move(device *i2c.PCA9685Driver) error {
	log.Printf("Servo Move Run Loop...\n")

	deg := 90

	// MIDDLE
	pulse := degree2pulse(deg)
	err := device.SetPWM(3, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(500 * time.Millisecond)

	// INC
	pulse = degree2pulse(deg + increase)
	err = device.SetPWM(3, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(1000 * time.Millisecond)

	// DEC
	pulse = degree2pulse(deg - increase)
	err = device.SetPWM(3, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}
	time.Sleep(1000 * time.Millisecond)

	// MIDDLE
	pulse = degree2pulse(deg)
	err = device.SetPWM(3, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}

	return err
}

func main() {

	// read some parameters from the command line
	if len(os.Args) > 0 {
		if len(os.Args) == 2 {
			// increase only
			increase, _ = strconv.Atoi(os.Args[1])
		}

		if len(os.Args) == 4 {
			// increase and servo pulse range
			increase, _ = strconv.Atoi(os.Args[1])
			servoMin, _ = strconv.Atoi(os.Args[2])
			servoMax, _ = strconv.Atoi(os.Args[3])
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
		gobot.Every(10*time.Second, func() {
			move(pca9685)
		})
	}

	robot := gobot.NewRobot("gobot",
		[]gobot.Connection{r},
		[]gobot.Device{pca9685},
		work,
	)

	robot.Start()
}
