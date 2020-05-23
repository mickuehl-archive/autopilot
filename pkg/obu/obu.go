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

		// Vehicle specific functions

		// Direction sets the steering direction (in deg)
		Direction(value int)
		// Throttle sets the speed (-100..0..100)
		Throttle(value int)
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
}
