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
	logger = log.New("module", "unittest")
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
		logger.Info("autopilot engaged")
		time.Sleep(5 * time.Second)
		logger.Info("autopilot done ...")
		ap.Stop() // stops the execution of this process, otherwise we have to wait for a SIGTERM
	}
	ap.AddWork(testdrive)

	// initialize the autopilot & vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically
}
