package main

import (
	"fmt"
	"os"
	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/parts"
	"time"
)

func main() {
	// create a virtual OBU
	obu := parts.NewVirtualOBU()
	// standard autopilot
	ap, err := autopilot.NewInstance(obu)
	if err != nil {
		fmt.Errorf("Error initializing the autopilot: %w", err)
		os.Exit(1)
	}
	defer ap.Shutdown()

	// initialize the vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// test the pilot
	time.Sleep(5 * time.Second)

	// cleanup should happen automatically
}
