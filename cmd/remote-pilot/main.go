package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/obu"
	"shadow-racer/autopilot/v1/pkg/parts"
)

var (
	logger log.Logger
	// the autopilot instance
	ap *autopilot.Autopilot
	// the obu in use
	ob obu.OnboardUnit
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

	// get command line options
	flag.StringVar(&obuType, "obu", "raspi", "Select an on-board unit implementation")
	flag.IntVar(&port, "port", 3000, "Port of the remote UI and API")
	flag.StringVar(&broker, "b", "tcp://localhost:1883", "MQTT Broker endpoint")
	flag.StringVar(&queue, "q", "shadow-racer/telemetry", "Default queue for telemetry data")
	flag.Parse()

	// create a OBU and autopilot
	if obuType == "virtual" {
		// create a virtual OBU for local testing
		ob = parts.NewVirtualOnboardUnit()
		// standard autopilot
		ap, _ = autopilot.NewInstance(ob)
	} else if obuType == "raspi" {
		// create a ShadowRacer OBU
		o, err := parts.NewRaspiOnboardUnit()
		if err != nil {
			os.Exit(1)
		}
		ob = o

		// standard autopilot
		ap, err = autopilot.NewInstance(ob)
		if err != nil {
			os.Exit(1)
		}
	} else {
		os.Exit(1) // should not happen
	}
	defer ap.Shutdown()

	//
	// add parts to the autopilot
	//

	// capture the camera stream
	ap.AddPart("camera", parts.NewLiveStreamCamera(fmt.Sprintf(":%d", port+1)))
	// maintain the vehicle state
	ap.AddPart("state", parts.NewVehicleState(ob))
	// "connected-car", send the vehicle state to a DC
	ap.AddPart("telemetry", parts.NewTelemetry(broker, queue))

	// start the http server as the remote pilot
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

	// initialize the autopilot and all parts
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically on exit
}
