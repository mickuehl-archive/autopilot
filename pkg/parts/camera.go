package parts

import (
	"os/exec"
)

type (
	// LiveStreamCamera encapsulates the PID to start/stop the camera process
	LiveStreamCamera struct {
		proc *exec.Cmd
	}
)

// NewLiveStreamCamera returns a new instance of a live stream camera
func NewLiveStreamCamera(port string) *LiveStreamCamera {
	c := &LiveStreamCamera{}
	return c
}

// Initialize prepares the camera server
func (c *LiveStreamCamera) Initialize() error {
	c.proc = exec.Command("nettop")

	err := c.proc.Start()
	return err
}

// Reset re-initializes the camera server
func (c *LiveStreamCamera) Reset() error {
	err := c.proc.Process.Kill()
	if err != nil {
		return err
	}
	return c.Initialize()
}

// Shutdown releases all camera server resources
func (c *LiveStreamCamera) Shutdown() error {
	return c.proc.Process.Kill()
}
