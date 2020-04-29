package parts

type (
	// StandardSpeedController represents a simple ESC
	StandardSpeedController struct {
		// -100 .. +100
		Throttle int
		// hardware config and values
		Cfg  ChannelCfg
		Data ChannelData
	}
)

// SetThrottle set the throttle
func (esc *StandardSpeedController) SetThrottle(value int) {
	logger.Debug("StandardSpeedController", "throttle", value)

	if value < -100 {
		value = -100
	} else if value > 100 {
		value = 100
	}
	esc.Throttle = value // write the sanitized value back

	if value == 0 {
		esc.Data.PulseOn = 0
		esc.Data.PulseOff = esc.Cfg.ZeroPulse
	} else {
		// ignore reversing for now
		esc.Data.PulseOn = esc.Cfg.BasePulse
		esc.Data.PulseOff = esc.Cfg.ZeroPulse + int(float32(esc.Cfg.MaxPulse-esc.Cfg.ZeroPulse)*float32(value)/100.0)
	}
}

// GetThrottle returns the current throttle value
func (esc *StandardSpeedController) GetThrottle() int {
	return esc.Throttle
}
