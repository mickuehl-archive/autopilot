package obu

type (
	// Config holds basic pre-sets
	Config struct {
		Frequency int // PCA9685 clock speed
		// actuators
		Steering *Channel
		Drive    *Channel
		Led1     *Channel // break
		Led2     *Channel
		Led3     *Channel // indicator, rear
		Led4     *Channel
		Led5     *Channel // indicator, front
		Led6     *Channel
	}
)
