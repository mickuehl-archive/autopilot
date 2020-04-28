package obu

import (
	log "github.com/majordomusio/log15"
)

type (
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
		PulseFunc     ChannelFunc // FIXME move this to the device driver
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
}

// Initialize s the OBU and all its components
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
