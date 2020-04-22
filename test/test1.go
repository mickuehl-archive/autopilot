package main

import (
	"fmt"
	"os"
	"shadow-racer/autopilot/v1/pkg/pilot"
	"time"
)

func main() {

	p, err := pilot.NewPilot()
	if err != nil {
		fmt.Errorf("Error initializing the pilot: %w", err)
		os.Exit(1)
	}
	defer p.Shutdown()
	p.Start()

	// test the servo
	p.Direction(0)
	time.Sleep(2 * time.Second)
	p.Direction(20)
	time.Sleep(2 * time.Second)
	p.Direction(-20)
	time.Sleep(2 * time.Second)
	p.Direction(0)
	time.Sleep(2 * time.Second)

	// cleanup should be done automatically
}
