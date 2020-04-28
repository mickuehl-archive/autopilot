package parts

/*
// StandardESCThrottle sets the motor speed [-100,+100]
func StandardESCThrottle(obu *obu.OnboardUnit, value int) {
	logger.Debug("StandardESCThrottle", "thr", value)

	ch := obu.Cfg.Drive

	if value < 0 {
		// assuming the ESC can not reverse
		value = 0
	} else if value > 100 {
		value = 100
	}

	pulseOn := 0
	pulseOff := 0
	if value == 0 {
		pulseOff = ch.ZeroPulse
	} else {
		// p.cfg.escZeroPulse + int(float32(p.cfg.escMaxPulse)*p.throttle)
		pulseOn = ch.ZeroPulse
		pulseOff = ch.ZeroPulse + int(float32(ch.MaxPulse-ch.ZeroPulse)*float32(value)/100.0)
	}

	// store the new values
	ch.CurMinPulse = pulseOn
	ch.CurMaxPulse = pulseOff

	// set the ESC pulse
	obu.PulseFunc(obu, ch.ChannelNo, pulseOn, pulseOff)
}
*/
