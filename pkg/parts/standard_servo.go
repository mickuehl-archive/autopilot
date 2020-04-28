package parts

/*
// StandardServoDirection sets the steering angle [-45,+45]
func StandardServoDirection(obu *obu.OnboardUnit, value int) {
	logger.Debug("StandardServoDirection", "deg", value)

	ch := obu.Cfg.Steering
	maxDegree := float32(ch.CustomValues[0])
	maxRange := ch.CustomValues[1]
	minRange := maxRange * -1
	trim := ch.CustomValues[2]

	pulseOn := ch.BasePulse
	if value < minRange {
		value = minRange
	} else if value > maxRange {
		value = maxRange
	}
	direction := (maxDegree / 2.0) + float32(value+trim)
	pulseOff := ch.MinPulse + int(float32(ch.MaxPulse-ch.MinPulse)/maxDegree*direction)

	// store the new values
	ch.CurMinPulse = pulseOn
	ch.CurMaxPulse = pulseOff

	// set the servo pulse
	obu.PulseFunc(obu, ch.ChannelNo, pulseOn, pulseOff)
}
*/
