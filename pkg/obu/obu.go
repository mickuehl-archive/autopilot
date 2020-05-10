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
		Mode     string  `json:"mode"` // some string
		Throttle float32 `json:"th"`   // 0 .. 100
		Steering float32 `json:"st"`   // in deg
		Heading  float32 `json:"head"` // heading of the vehicle 0 -> North, 90 -> East ...
		TS       int64   `json:"ts"`   // timestamp
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
}
