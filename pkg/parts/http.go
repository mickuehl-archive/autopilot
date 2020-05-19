package parts

import (
	"encoding/json"
	"net/http"
	"shadow-racer/autopilot/v1/pkg/eventbus"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type (
	// RemoteState is the current driving state info
	RemoteState struct {
		Mode      string  `json:"mode"`
		Recording bool    `json:"recording"`
		Throttle  float32 `json:"th"`
		Steering  float32 `json:"st"`
		TS        int64   `json:"ts"`
	}
)

// StartHTTPServer launches a http server that serves static content from ./public and exposes a websocket endpoint
func StartHTTPServer(addr string) error {
	// most basic HTTP server in golang
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/ws", http.HandlerFunc(wsHandler))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		logger.Error("Error", "msg", err.Error())
	}

	return err
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
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
					var state RemoteState
					err := json.Unmarshal(msg, &state)
					if err == nil {
						eventbus.InstanceOf().Publish("rc/state", state)
					}
				}
			}
		}()

		go func() {
			// sender
			ch := eventbus.InstanceOf().Subscribe("state/vehicle")
			for {
				vehicle := <-ch
				data, err := json.Marshal(&vehicle)
				if err == nil {
					wsutil.WriteServerMessage(conn, opcode, data)
				}
			}
		}()
	}
}
