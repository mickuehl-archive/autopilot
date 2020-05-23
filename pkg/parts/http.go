package parts

import (
	"encoding/json"
	"net/http"
	"shadow-racer/autopilot/v1/pkg/eventbus"
	"shadow-racer/autopilot/v1/pkg/metrics"

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
	// collect metrics
	metrics.NewMeter(mHUDReceive)
	metrics.NewMeter(mHUDUpdate)
	metrics.NewMeter(mImageReceive)

	// most basic HTTP server in golang
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/state", http.HandlerFunc(wsStateHandler))
	http.Handle("/image", http.HandlerFunc(wsImageHandler))

	return http.ListenAndServe(addr, nil)
}

func wsStateHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error("Error upgrading HTTP to WS", "err", err.Error())
		return
	}

	var opcode ws.OpCode

	go func() {
		// receiver
		defer conn.Close()
		for {
			msg, op, err := wsutil.ReadClientData(conn)
			opcode = op

			if err != nil {
				logger.Error("Error receiving WS message", "err", err.Error())
				break // FIXME abort on the first error, really ?
			} else {
				var state RemoteState
				err := json.Unmarshal(msg, &state)
				if err == nil {
					eventbus.InstanceOf().Publish(topicRCStateReceive, state)
				}
			}
			metrics.Mark(mHUDReceive)
		}
	}()

	go func() {
		// sender
		ch := eventbus.InstanceOf().Subscribe(topicRCStateUpdate)
		for {
			vehicle := <-ch
			data, err := json.Marshal(&vehicle)
			if err == nil {
				wsutil.WriteServerMessage(conn, opcode, data)
				metrics.Mark(mHUDUpdate)
			}
		}
	}()
}

func wsImageHandler(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		logger.Error("Error upgrading HTTP to WS", "err", err.Error())
		return
	}

	go func() {
		defer conn.Close()
		for {
			msg, _, err := wsutil.ReadClientData(conn)

			if err != nil {
				logger.Error("Error receiving WS message", "err", err.Error())
				break // FIXME abort on the first error, really ?
			} else {
				eventbus.InstanceOf().Publish(topicImageReceive, msg) // FIXME check if we can pass a pointer to msg
			}
			metrics.Mark(mImageReceive)
		}
	}()
}
