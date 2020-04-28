package parts

import (
	"shadow-racer/autopilot/v1/pkg/pilot"
)

// NewVirtualOBU creates a software-only OBU configration for local unit testing
func NewVirtualOBU() *pilot.OnboardUnit {
	cfg := &pilot.Config{
		Frequency: 50,
		Steering: &pilot.Channel{
			ChannelNo:    3,
			MinPulse:     180, // values from a real servo (MG996R)
			MaxPulse:     590,
			BasePulse:    100,
			ZeroPulse:    385,
			InitPulse:    -1,                // not used
			CustomValues: []int{180, 30, 0}, // Max Range, Range Limit, trim
		},
		Drive: &pilot.Channel{
			ChannelNo: 0,
			MinPulse:  100, // not sure about these values
			MaxPulse:  100,
			BasePulse: 1000,
			ZeroPulse: 1300,
			InitPulse: 2000,
		},
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
