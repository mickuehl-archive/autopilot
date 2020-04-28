package obu

type (
	// Channel holds the configuration for controlling an actuator (e.g. servo or ESC)
	Channel struct {
		// channel number
		ChannelNo int
		// min pulse length out of 4096
		MinPulse int
		// max pulse length out of 4096
		MaxPulse int
		// base puls to start from (min value)
		BasePulse int
		// zero pulse where the actor (e.g. servo) is in a neutral position
		ZeroPulse int
		// init pulse used to initialize/reset the channel. Set to -1 if not used
		InitPulse int
		// CustomValues allows to pass additional config parameters, e.g. range limits or trim values
		CustomValues []int
		// current min/max pulse settings
		CurMinPulse int
		CurMaxPulse int
	}
)
