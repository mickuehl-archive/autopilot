package parts

import (
	"shadow-racer/autopilot/v1/pkg/pilot"
)

// NewVirtualOBU creates a software-only OBU configration for local unit testing
func NewVirtualOBU() *pilot.OnboardUnit {
	cfg := &pilot.Config{
		Frequency: 50,
		Steering:  newVirtualChannel(3),
		Drive:     newVirtualChannel(0),
	}

	obu := &pilot.OnboardUnit{
		Cfg:           cfg,
		InitFunc:      VirtualOBUInitialize,
		ShutdownFunc:  VirtualOBUShutdown,
		DirectionFunc: StandardServoDirection,
		ThrottleFunc:  StandardESCThrottle,
		PulseFunc:     VirtualOBUPulse,
	}
	return obu
}

func newVirtualChannel(n int) *pilot.Channel {
	return &pilot.Channel{
		ChannelNo: n,
		MinPulse:  100,
		MaxPulse:  500,
		BasePulse: 100,
		ZeroPulse: 300,
		InitPulse: 2000,
	}
}

// VirtualOBUInitialize s the pilot and all its components
func VirtualOBUInitialize(cfg *pilot.Config) error {
	logger.Debug("VirtualOBUInitialize")
	return nil
}

// VirtualOBUShutdown stops & resets all components
func VirtualOBUShutdown(cfg *pilot.Config) error {
	logger.Debug("VirtualOBUShutdown")
	return nil
}

// VirtualOBUPulse sets the pulse calues of a channel
func VirtualOBUPulse(obu *pilot.OnboardUnit, ch, min, max int) {
	logger.Debug("VirtualOBUPulse", "channel", ch, "min_pulse", min, "max_pulse", max)
}
