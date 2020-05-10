package parts

type (
	// StandardServo represents a servo
	StandardServo struct {
		// servo configuration
		MaxDegree int
		MaxRange  int
		MinRange  int
		Trim      int
		// servo state
		Direction int
		// hardware config and values
		Cfg ChannelCfg
	}
)

// SetAngle sets the servo angle
func (s *StandardServo) SetAngle(value int) (int, int) {
	//logger.Debug("StandardServo", "deg", value)

	on := s.Cfg.BasePulse
	off := 0
	if value < s.MinRange {
		value = s.MinRange
	} else if value > s.MaxRange {
		value = s.MaxRange
	}
	s.Direction = value // write the sanitized value back

	// calculate the OFF value
	deg := (float32(s.MaxDegree) / 2.0) + float32(value+s.Trim)
	off = s.Cfg.MinPulse + int(float32(s.Cfg.MaxPulse-s.Cfg.MinPulse)/float32(s.MaxDegree)*deg)

	return on, off
}

// GetAngle returns the current steering direction
func (s *StandardServo) GetAngle() int {
	return s.Direction
}
