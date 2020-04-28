package parts

type (
	VirtualOnboardUnit struct {
		//
	}
)

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
