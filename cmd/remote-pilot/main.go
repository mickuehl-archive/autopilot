package main

import (
	"encoding/json"

	log "github.com/majordomusio/log15"

	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type (
	// State is the current driving state info
	State struct {
		Throttle float32 `json:"th,omitempty"`
		Steering float32 `json:"st,omitempty"`
		Mode     string  `json:"mode,omitempty"`
		TS       int64   `json:"ts,omitempty"`
	}

	// HUD holds the state info and additional data to be shown on the display
	HUD struct {
		Throttle float32 // 0 .. 1.0
		Steering float32 // in deg
		Heading  float32 // heading of the car
		TS       int64
	}
)

var (
	logger log.Logger
)

func main() {
	logger = log.New("module", "main")

	// most basic HTTP server in golang
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/hud", http.HandlerFunc(hudWSHandler))

	logger.Info("Listening on :3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		logger.Error("%v", err)
	}
}

func hudWSHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err == nil {
		defer conn.Close()

		go func() {
			// receiver
			for {
				msg, _, err := wsutil.ReadClientData(conn)

				if err != nil {
					break // FIXME abort on the first error, really ?
				} else {
					var state State
					err := json.Unmarshal(msg, &state)
					if err == nil {
						logger.Debug("Rcvd", "state", state)
					}
				}

				//err = wsutil.WriteServerMessage(conn, op, msg)
				//if err != nil {
				//	// handle error
				//}
			}
		}()
	}
}
