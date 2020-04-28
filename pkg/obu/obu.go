package obu

import (
	log "github.com/majordomusio/log15"
)

type (
	// OnboardUnit provides an interface to actual harware
	OnboardUnit interface {
		// Initialize prepares the device
		Initialize() error
		// Reset re-initializes the device
		Reset() error
		// Shutdown releases all resources
		Shutdown() error
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
}
