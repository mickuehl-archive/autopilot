package main

import (
	"fmt"
	"os"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/parts"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "remote-pilot")
}

func main() {

	// create a virtual OBU
	obu := parts.NewVirtualOnboardUnit()
	// standard autopilot
	ap, err := autopilot.NewInstance(obu)
	if err != nil {
		fmt.Errorf("Error initializing the autopilot: %w", err)
		os.Exit(1)
	}
	defer ap.Shutdown()

	// add a VERY simplistic autopilot activity
	testdrive := func() {
		logger.Info("RemotePilot engaged")

		err := parts.StartHTTPServer(":3000") // FIXME configuration
		if err != nil {
			logger.Error("RemotePilot aborted")
			ap.Stop()
		}

		logger.Info("RemotePilot disengaged")
	}
	ap.AddWork(testdrive)

	// initialize the autopilot & vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically
}
