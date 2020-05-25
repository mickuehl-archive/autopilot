package parts

import (
	"os/exec"
	"shadow-racer/autopilot/v1/pkg/metrics"
	"shadow-racer/autopilot/v1/pkg/sharedm"
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
	metrics.NewMeter(mImageReceive)

	c.proc = exec.Command("./camera.py") // FIXME pass parameters
	stdout, err := c.proc.StdoutPipe()
	if err != nil {
		logger.Error("Error getting STDOUT", "err", err.Error())
		return err
	}

	// read images from STDOUT and update the shared memory
	go func() {
		var buffer []byte
		buffer = make([]byte, 100000) // FIXME depending on the camera resolution !!

		for {
			n, err := stdout.Read(buffer)
			if err == nil {
				if n > 0 {
					sharedm.StoreBytes(memImageRaw, buffer[:n])
					metrics.Mark(mImageReceive)
				}
			}
		}
	}()

	return c.proc.Start()
}

// FIXME check if process is still running !

// Reset re-initializes the camera server
func (c *LiveStreamCamera) Reset() error {
	if c.proc.Process != nil {
		err := c.proc.Process.Kill()
		if err != nil {
			return err
		}
	}
	return c.Initialize()
}

// Shutdown releases all camera server resources
func (c *LiveStreamCamera) Shutdown() error {
	if c.proc.Process != nil {
		return c.proc.Process.Kill()
	}
	return nil
}
