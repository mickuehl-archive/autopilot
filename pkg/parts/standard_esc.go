package parts

type (
	// StandardSpeedController represents a simple ESC
	StandardSpeedController struct {
		// -100 .. +100
		Throttle   int
		LowerLimit int
		UpperLimit int
		// hardware config and values
		Cfg ChannelCfg
	}
)

// SetThrottle set the throttle
func (esc *StandardSpeedController) SetThrottle(value int) (int, int) {

	if value < esc.LowerLimit {
		value = esc.LowerLimit
	} else if value > esc.UpperLimit {
		value = esc.UpperLimit
	}
	esc.Throttle = value // write the sanitized value back

	on := 0
	off := 0

	if value == 0 {
		off = esc.Cfg.ZeroPulse
	} else {
		// ignore reversing for now
		on = esc.Cfg.BasePulse
		off = esc.Cfg.ZeroPulse + int(float32(esc.Cfg.MaxPulse-esc.Cfg.ZeroPulse)*float32(value)/100.0)
	}

	return on, off
}

// GetThrottle returns the current throttle value
func (esc *StandardSpeedController) GetThrottle() int {
	return esc.Throttle
}
