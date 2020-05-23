package parts

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
	obu.running = true
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

// Direction sets the steering direction (in deg)
func (obu *VirtualOnboardUnit) Direction(value int) {
}

// Throttle sets the speed (-100..0..100)
func (obu *VirtualOnboardUnit) Throttle(value int) {
}
