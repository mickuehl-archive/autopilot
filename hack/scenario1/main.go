package main

import (
	"fmt"
	"os"
	"strconv"
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
	// defaults
	speed := 20
	t := 1500 // ms

	// read values from the command line
	if len(os.Args) == 3 {
		speed, _ = strconv.Atoi(os.Args[1])
		t, _ = strconv.Atoi(os.Args[1])
	}

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
		logger.Info("Start the driving scenario")

		// activate the taillights, flashing
		obu.TailLights(4000, true)

		// simple programm
		obu.Direction(0)
		obu.Throttle(0)
		time.Sleep(1 * time.Second)

		obu.Throttle(speed)
		time.Sleep(time.Duration(t) * time.Millisecond)

		obu.Throttle(0)
		obu.TailLightsOff()

		logger.Info("Done ...")
		ap.Stop()
	}
	ap.AddWork(testdrive)

	// initialize the autopilot
	ap.Initialize()

	// activate the autopilot
	ap.Activate()
}
