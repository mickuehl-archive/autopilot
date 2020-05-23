package parts

const (
	topicRCStateReceive = "rc/state"
	topicRCStateSend    = "state/vehicle"
)

// NewBMS390DMH creates an instance of a BMS390DMH servo
func NewBMS390DMH(n int) *StandardServo {
	return &StandardServo{
		MaxDegree: 60,
		MaxRange:  30,
		MinRange:  -30,
		Trim:      3,
		Direction: 0,
		Cfg: ChannelCfg{
			N:         n,
			MinPulse:  1230,
			MaxPulse:  1350,
			BasePulse: 1000,
			ZeroPulse: 1290,
			InitPulse: -1,
		},
	}
}

// NewMG996R creates an instance of a MG996R servo
func NewMG996R(n int) *StandardServo {
	return &StandardServo{
		MaxDegree: 180,
		MaxRange:  30,
		MinRange:  -30,
		Trim:      0,
		Direction: 0,
		Cfg: ChannelCfg{
			N:         n,
			MinPulse:  180,
			MaxPulse:  590,
			BasePulse: 100,
			ZeroPulse: 385,
			InitPulse: -1,
		},
	}
}

// NewWP40 creates an instance of a Reely WP40 speed controller
func NewWP40(n int) *StandardSpeedController {
	return &StandardSpeedController{
		Throttle: 0,
		Limit:    100,
		Cfg: ChannelCfg{
			N:         n,
			MinPulse:  1000, // not sure
			MaxPulse:  1400, // not sure
			BasePulse: 1000,
			ZeroPulse: 1300,
			InitPulse: 2000,
		},
	}
}
