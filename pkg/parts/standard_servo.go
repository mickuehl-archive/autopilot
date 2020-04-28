package parts

import (
	"shadow-racer/autopilot/v1/pkg/pilot"
)

const (
	// MaxDegree is the maximum movement of the servo
	MaxDegree = 180.0
	// MaxRange is the allowed degree of freedom on the positive/right side. The max. value is MaxDegree / 2
	MaxRange = 30
	// MinRange is the allowed degree of freedom on the negative/left side. The max. value is MaxDegree / 2
	MinRange = -30
)

// StandardServoDirection sets the steering angle [-45,+45]
func StandardServoDirection(obu *pilot.OnboardUnit, value int) {
	logger.Debug("StandardServoDirection", "deg", value)

	ch := obu.Cfg.Steering

	if value < MinRange {
		value = MinRange
	} else if value > MaxRange {
		value = MaxRange
	}
	direction := (MaxDegree / 2) + value
	pulse := ch.BasePulse + int(float32(ch.MaxPulse-ch.MinPulse)/MaxDegree*float32(direction))

	// set the servo pulse
	obu.PulseFunc(obu, ch.ChannelNo, ch.BasePulse, pulse)
}
