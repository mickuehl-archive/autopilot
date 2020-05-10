package parts

import (
	"shadow-racer/autopilot/v1/pkg/eventbus"
	"shadow-racer/autopilot/v1/pkg/obu"

	"github.com/majordomusio/commons/pkg/util"
)

const (
	// ServoRange is the max allowed servo movement
	ServoRange = 30.0
)

type (
	// VirtualOnboardUnit holds data needed to simulate an OBU
	VirtualOnboardUnit struct {
		running bool
	}
)

// NewVirtualOnboardUnit creates a new instance of a virtual OBU, for e.g. unit testing of the framework
func NewVirtualOnboardUnit() *VirtualOnboardUnit {
	obu := &VirtualOnboardUnit{
		running: false,
	}
	return obu
}

// Initialize prepares the device
func (obu *VirtualOnboardUnit) Initialize() error {
	logger.Debug("initialize the obu")

	obu.running = true
	go remoteStateHandler()

	return nil
}

// Reset re-initializes the device
func (obu *VirtualOnboardUnit) Reset() error {
	logger.Debug("reset the obu")

	return nil
}

// Shutdown releases all resources
func (obu *VirtualOnboardUnit) Shutdown() error {
	logger.Debug("shutdown the obu")

	obu.running = false

	return nil
}

func remoteStateHandler() {
	ch := eventbus.InstanceOf().Subscribe("rc/state")
	for {
		evt := <-ch
		state := evt.Data.(RemoteState)

		vehicle := obu.Vehicle{
			Mode:     state.Mode,
			Steering: 100 * ((ServoRange / 90.0) * state.Steering),
			Throttle: 100 * state.Throttle,
			Heading:  360,
			TS:       util.Timestamp(),
		}

		eventbus.InstanceOf().Publish("rc/vehicle", vehicle)
	}
}
