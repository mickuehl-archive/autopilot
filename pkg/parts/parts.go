package parts

import (
	log "github.com/majordomusio/log15"
)

type (
	// Part provides an interface to either a physical or logical part of the vehicle
	Part interface {
		// Initialize prepares the device
		Initialize() error
		// Reset re-initializes the device
		Reset() error
		// Shutdown releases all resources
		Shutdown() error
	}

	// ChannelFunc sets the pulse values of a channel
	ChannelFunc func(int, int, int)

	// ChannelCfg holds the configuration for controlling an actuator e.g. a servo or ESC
	ChannelCfg struct {
		// channel number
		N int
		// min pulse length out of 4096
		MinPulse int
		// max pulse length out of 4096
		MaxPulse int
		// base puls to start from (min value)
		BasePulse int
		// zero pulse where the actor (e.g. servo) is in a neutral position
		ZeroPulse int
		// init pulse used to initialize/reset the channel. Set to -1 if not used
		InitPulse int
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "parts")
}
