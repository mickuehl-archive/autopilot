package parts

import (
	"fmt"
	"shadow-racer/autopilot/v1/pkg/eventbus"
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
	go steeringHandler()

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

func steeringHandler() {
	ch := eventbus.InstanceOf().Subscribe("obu/steering")
	for {
		de := <-ch
		fmt.Printf("obu/steering %v\n", de)
	}
}
