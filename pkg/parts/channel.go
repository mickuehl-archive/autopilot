package parts

type (
	// ChannelFunc sets the pulse values of a channel
	ChannelFunc func(int, int, int)

	// ChannelCfg holds the configuration for controlling an actuator (e.g. servo or ESC)
	ChannelCfg struct {
		// channel number
		N int
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
	}

	// ChannelData represents the current channel settings
	ChannelData struct {
		// channel number
		N int
		// current pulse settings
		PulseOn  int
		PulseOff int
	}
)
