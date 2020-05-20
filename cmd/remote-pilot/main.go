package main

import (
	"flag"
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

	// command line parameters
	var obuType string
	var port int
	var broker string
	var queue string

	// the autopilot instance
	var ap *autopilot.Autopilot

	// get command line options
	flag.StringVar(&obuType, "obu", "raspi", "Select an on-board unit implementation")
	flag.IntVar(&port, "port", 3000, "Port of the remote UI and API")
	flag.StringVar(&broker, "b", "tcp://localhost:1883", "MQTT Broker endpoint")
	flag.StringVar(&queue, "q", "shadow-racer/telemetry", "Default queue for telemetry data")

	flag.Parse()

	// create a OBU and add it to an autopilot
	if obuType == "raspi" {
		// create a ShadowRacer OBU
		obu, err := parts.NewRaspiOnboardUnit()
		if err != nil {
			os.Exit(1)
		}

		// standard autopilot
		ap, err = autopilot.NewInstance(obu)
		if err != nil {
			os.Exit(1)
		}
	} else if obuType == "virtual" {
		// create a virtual OBU for local testing
		obu := parts.NewVirtualOnboardUnit()
		// standard autopilot
		ap, _ = autopilot.NewInstance(obu)
	} else {
		os.Exit(1)
	}
	defer ap.Shutdown()

	// add parts to the autopilot
	ap.AddPart("camera", parts.NewLiveStreamCamera(fmt.Sprintf(":%d", port+1)))
	ap.AddPart("telemetry", parts.NewTelemetry(broker, queue))

	// add a http server as the remote pilot
	remotepilot := func() {
		logger.Info("RemotePilot engaged")

		err := parts.StartHTTPServer(fmt.Sprintf(":%d", port))
		if err != nil {
			logger.Error("Error starting the HTTP listener", "err", err.Error())
			ap.Stop()
		}

		logger.Info("RemotePilot disengaged")
	}
	ap.AddWork(remotepilot)

	// initialize the autopilot & vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically
}
