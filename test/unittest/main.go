package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/eventbus"
	"shadow-racer/autopilot/v1/pkg/parts"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "unittest")
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

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
		logger.Info("Autopilot engaged")

		for i := 0; i < 10; i++ {
			eventbus.InstanceOf().Publish("obu/steering", r.Intn(30))
		}

		time.Sleep(10 * time.Second)
		logger.Info("Autopilot done ...")
		ap.Stop()
	}
	ap.AddWork(testdrive)

	// initialize the autopilot & vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically
}
