package main

import (
	"fmt"
	"os"
	"shadow-racer/autopilot/v1/pkg/pilot"
	"time"
)

func servo(p *pilot.Pilot) {
	p.Direction(0)
	time.Sleep(2 * time.Second)
	p.Direction(20)
	time.Sleep(2 * time.Second)
	p.Direction(-20)
	time.Sleep(2 * time.Second)
	p.Direction(0)
	time.Sleep(2 * time.Second)
}

func drive(p *pilot.Pilot) {
	throttle := 0.15
	fmt.Printf("Setting throttle to %f\n", throttle)
	p.Throttle(float32(throttle))
	time.Sleep(2 * time.Second)

	throttle = 0.2
	fmt.Printf("Setting throttle to %f\n", throttle)
	p.Throttle(float32(throttle))
	time.Sleep(2 * time.Second)

	throttle = 0.25
	fmt.Printf("Setting throttle to %f\n", throttle)
	p.Throttle(float32(throttle))
	time.Sleep(2 * time.Second)

	throttle = 0.1
	fmt.Printf("Setting throttle to %f\n", throttle)
	p.Throttle(float32(throttle))
	time.Sleep(3 * time.Second)

}

func main() {
	p, err := pilot.NewPilot()
	if err != nil {
		fmt.Errorf("Error initializing the pilot: %w", err)
		os.Exit(1)
	}
	defer p.Shutdown()
	p.Start()

	// test the pilot
	// go servo(p)

	// test the esc
	drive(p)

	// cleanup should be done automatically
}
