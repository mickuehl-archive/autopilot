package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/eventbus"
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
		logger.Info("Remote-Pilot engaged")

		// most basic HTTP server in golang
		http.Handle("/", http.FileServer(http.Dir("./public")))
		http.Handle("/hud", http.HandlerFunc(hudWSHandler))

		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			logger.Error("Remote-Pilot aborted ...")
			ap.Stop()
		}

		/*
			for i := 0; i < 10; i++ {
				eventbus.InstanceOf().Publish("obu/steering", r.Intn(30))
			}
		*/

		logger.Info("Remote-Pilot done ...")
	}
	ap.AddWork(testdrive)

	// initialize the autopilot & vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// cleanup should happen automatically
}

func hudWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err == nil {

		var opcode ws.OpCode

		go func() {
			// receiver
			defer conn.Close()
			for {
				msg, op, err := wsutil.ReadClientData(conn)
				opcode = op

				if err != nil {
					break // FIXME abort on the first error, really ?
				} else {
					var state parts.RemoteState
					err := json.Unmarshal(msg, &state)
					if err == nil {
						eventbus.InstanceOf().Publish("remote/state", state)

						//logger.Debug("Rcvd", "state", state)
					}
				}

				//err = wsutil.WriteServerMessage(conn, op, msg)
				//if err != nil {
				//	// handle error
				//}
			}
		}()

		go func() {
			// sender
			ch := eventbus.InstanceOf().Subscribe("remote/hud")
			for {
				hud := <-ch
				data, err := json.Marshal(&hud)
				if err == nil {
					//fmt.Printf("remote/hud %s\n", string(data))
					wsutil.WriteServerMessage(conn, opcode, data)
				}
			}
		}()
	}
}
