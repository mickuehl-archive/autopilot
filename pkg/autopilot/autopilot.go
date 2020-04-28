package autopilot

import (
	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/obu"
)

type (
	// Autopilot holds all resources needed to pilot a vehicle
	Autopilot struct {
		obu *obu.OnboardUnit
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "autopilot")
}

// NewInstance creates and initializes a new autopilot instance
func NewInstance(obu *obu.OnboardUnit) (*Autopilot, error) {
	ap := &Autopilot{
		obu: obu,
	}
	return ap, nil
}

// Initialize prepares the autopilot instance
func (ap *Autopilot) Initialize() error {
	logger.Debug("initialize")
	err := ap.obu.Initialize()
	if err != nil {
		// FIXME something else?
		return err
	}
	// FIXME do autopilot stuff here e.g
	ap.obu.Direction(0)
	ap.obu.Throttle(0)

	return err
}

// Activate will start the autopilot
func (ap *Autopilot) Activate() error {
	logger.Debug("activate")
	// FIXME do autopilot stuff here
	return nil
}

// Stop cancels the autopilot and stops the vehicle
func (ap *Autopilot) Stop() error {
	logger.Debug("stop")
	// FIXME do autopilot stuff here
	return nil
}

// Shutdown finalizes the autopilot and releases all resources
func (ap *Autopilot) Shutdown() error {
	logger.Debug("shutdown")
	err := ap.obu.Shutdown()
	if err != nil {
		// FIXME something else?
		return err
	}
	// FIXME do autopilot stuff here
	return err
}
