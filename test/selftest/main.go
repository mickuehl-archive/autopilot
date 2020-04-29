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

		obu.TailLights(4000, true)
		time.Sleep(5 * time.Second)
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
