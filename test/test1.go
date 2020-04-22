package main

import (
	"fmt"
	"os"
	"shadow-racer/autopilot/v1/pkg/pilot"
	"strconv"
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
	throttle := 0.0
	if len(os.Args) > 1 {
		i, _ := strconv.Atoi(os.Args[1])
		throttle = float64(i) / 100.0
	}

	fmt.Printf("Setting throttle to %f\n", throttle)
	p.Throttle(float32(throttle))
	time.Sleep(2 * time.Second)
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
