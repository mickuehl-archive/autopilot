package obu

import (
	log "github.com/majordomusio/log15"
)

type (
	// OnboardUnit provides an interface to actual hardware
	OnboardUnit interface {
		// Initialize prepares the device
		Initialize() error
		// Reset re-initializes the device
		Reset() error
		// Shutdown releases all resources
		Shutdown() error
	}

	// Vehicle holds state information of a generic vehicle
	Vehicle struct {
		Mode        string  `json:"mode"` // some string, e.g. DRIVING, STOPPED, etc
		Throttle    float32 `json:"th"`   // -100 .. 100
		Steering    float32 `json:"st"`   // in deg, 0 is straight ahead
		Heading     float32 `json:"head"` // heading of the vehicle 0 -> North, 90 -> East ...
		Recording   bool    `json:"recording"`
		RecordingTS int64   `json:"recording_ts"`
		TS          int64   `json:"ts"` // timestamp
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
}

// Clone returns a deep copy the vehicle state
func (v *Vehicle) Clone() *Vehicle {
	return &Vehicle{
		Mode:        v.Mode,
		Throttle:    v.Throttle,
		Steering:    v.Steering,
		Heading:     v.Heading,
		Recording:   v.Recording,
		RecordingTS: v.RecordingTS,
		TS:          v.TS,
	}
}
