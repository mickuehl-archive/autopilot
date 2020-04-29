package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/parts"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "selftest")
}

func main() {
	// create a virtual OBU
	obu, err := parts.NewRaspiOnboardUnit()
	if err != nil {
		fmt.Errorf("Error initializing the OBU: %w", err)
		os.Exit(1)
	}
	// standard autopilot
	ap, err := autopilot.NewInstance(obu)
	if err != nil {
		fmt.Errorf("Error initializing the autopilot: %w", err)
		os.Exit(1)
	}
	defer ap.Shutdown()

	testdrive := func() {
		logger.Info("Starting the component self-test")

		// activate the taillights, flashing
		obu.TailLights(4000, true)

		// test the servo
		obu.Direction(0)
		time.Sleep(1 * time.Second)
		obu.Direction(30)
		time.Sleep(1 * time.Second)
		obu.Direction(-30)
		time.Sleep(1 * time.Second)
		obu.Direction(0)

		// test the ESC
		obu.Throttle(20)
		time.Sleep(1 * time.Second)
		obu.Throttle(25)
		time.Sleep(1 * time.Second)
		obu.Throttle(0)
		time.Sleep(1 * time.Second)

		//time.Sleep(5 * time.Second)
		obu.TailLightsOff()

		logger.Info("Selftest done ...")
		ap.Stop()
	}
	ap.AddWork(testdrive)

	// initialize the autopilot
	ap.Initialize()

	// activate the autopilot
	ap.Activate()
}
