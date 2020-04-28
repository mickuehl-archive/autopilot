package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/autopilot"
	"shadow-racer/autopilot/v1/pkg/obu"
	"shadow-racer/autopilot/v1/pkg/parts"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "selftest")
}

// NewVirtualOBU creates a software-only OBU configration for local unit testing
func NewVirtualOBU() *obu.OnboardUnit {
	cfg := &obu.Config{
		Frequency: 50,
		Steering: &obu.Channel{
			ChannelNo:    3,
			MinPulse:     180, // values from a real servo (MG996R)
			MaxPulse:     590,
			BasePulse:    100,
			ZeroPulse:    385,
			InitPulse:    -1,                // not used
			CustomValues: []int{180, 30, 0}, // Max Range, Range Limit, trim
		},
		Drive: &obu.Channel{
			ChannelNo: 0,
			MinPulse:  100, // not sure about these values
			MaxPulse:  100,
			BasePulse: 1000,
			ZeroPulse: 1300,
			InitPulse: 2000,
		},
	}

	obu := &obu.OnboardUnit{
		Cfg:           cfg,
		InitFunc:      VirtualOBUInitialize,
		ShutdownFunc:  VirtualOBUShutdown,
		DirectionFunc: parts.StandardServoDirection,
		ThrottleFunc:  parts.StandardESCThrottle,
		PulseFunc:     VirtualOBUPulse,
	}
	return obu
}

// VirtualOBUInitialize s the pilot and all its components
func VirtualOBUInitialize(cfg *obu.Config) error {
	logger.Debug("VirtualOBUInitialize")
	return nil
}

// VirtualOBUShutdown stops & resets all components
func VirtualOBUShutdown(cfg *obu.Config) error {
	logger.Debug("VirtualOBUShutdown")
	return nil
}

// VirtualOBUPulse sets the pulse calues of a channel
func VirtualOBUPulse(obu *obu.OnboardUnit, ch, min, max int) {
	logger.Debug("VirtualOBUPulse", "channel", ch, "min_pulse", min, "max_pulse", max)
}

func main() {
	// create a virtual OBU
	obu := NewVirtualOBU()
	// standard autopilot
	ap, err := autopilot.NewInstance(obu)
	if err != nil {
		fmt.Errorf("Error initializing the autopilot: %w", err)
		os.Exit(1)
	}
	defer ap.Shutdown()

	// initialize the vehicle
	ap.Initialize()

	// activate the autopilot
	ap.Activate()

	// test the pilot
	time.Sleep(5 * time.Second)

	// cleanup should happen automatically
}
