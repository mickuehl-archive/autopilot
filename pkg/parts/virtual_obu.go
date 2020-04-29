package parts

type (
	// VirtualOnboardUnit holds data needed to simulate an OBU
	VirtualOnboardUnit struct {
		// nothing yet
	}
)

// NewVirtualOnboardUnit creates a new instance of a virtual OBU, for e.g. unit testing of the framework
func NewVirtualOnboardUnit() *VirtualOnboardUnit {
	obu := &VirtualOnboardUnit{}
	return obu
}

// Initialize prepares the device
func (obu *VirtualOnboardUnit) Initialize() error {
	logger.Debug("initialize")
	return nil
}

// Reset re-initializes the device
func (obu *VirtualOnboardUnit) Reset() error {
	logger.Debug("reset")
	return nil
}

// Shutdown releases all resources
func (obu *VirtualOnboardUnit) Shutdown() error {
	logger.Debug("shutdown")
	return nil
}
