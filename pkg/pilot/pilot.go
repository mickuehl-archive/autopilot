package pilot

import (
	log "github.com/majordomusio/log15"
)

type (
	// Channel holds the configuration for controlling
	Channel struct {
		// channel number
		n int
		// min pulse length out of 4096
		minPulse int
		// max pulse length out of 4096
		maxPulse int
		// base puls to start from (min value)
		basePulse int
		// zero pulse where the actor (e.g. servo) is in a neutral position
		zeroPulse int
		// init pulse used to initialize/reset the channel
		initPulse int
		// current min/max pulse settings
		curMinPulse int
		curMaxPulse int
	}

	// Config holds basic pre-sets
	Config struct {
		Frequency int // PCA9685 clock speed
		// actuators
		Servo Channel
		Esc   Channel
		Led1  Channel
		Led2  Channel
	}

	// StateFunc is a simple interface for e.g. the init function
	StateFunc func(*Config) error
	// ActuatorFunc controlls a channel
	ActuatorFunc func(*OnboardUnit, int)
	// ChannelFunc sets the pulse calues of a channel
	ChannelFunc func(*OnboardUnit, int, int, int)

	// OnboardUnit is an abstraction of the hardware controlling the vehicle
	OnboardUnit struct {
		Cfg           *Config
		InitFunc      StateFunc
		ShutdownFunc  StateFunc
		DirectionFunc ActuatorFunc
		ThrottleFunc  ActuatorFunc
		PulseFunc     ChannelFunc
	}

	// Pilot dsdskjf dsfhdshf
	Pilot interface {
		// Initializes the pilot and all its components
		Initialize() error
		// Shutdown stops & resets all components
		Shutdown() error
		// Direction sets the steering angle [-45,+45]
		Direction(value int)
		// Throttle sets the motor speed [-100,+100]
		Throttle(value int)
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "pilot")
}

// Initialize s the pilot and all its components
func (obu *OnboardUnit) Initialize() error {
	logger.Debug("initialize")
	return obu.InitFunc(obu.Cfg)
}

// Shutdown stops & resets all components
func (obu *OnboardUnit) Shutdown() error {
	logger.Debug("shutdown")
	return obu.ShutdownFunc(obu.Cfg)
}

// Direction sets the steering angle [-45,+45]
func (obu *OnboardUnit) Direction(value int) {
	logger.Debug("direction", "deg", value)
	obu.DirectionFunc(obu, value)
}

// Throttle sets the motor speed [-100,+100]
func (obu *OnboardUnit) Throttle(value int) {
	logger.Debug("throttle", "thr", value)
	obu.ThrottleFunc(obu, value)
}
