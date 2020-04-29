package autopilot

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/obu"
)

type (
	// Autopilot holds all resources needed to pilot a vehicle
	Autopilot struct {
		// an instance of the device we are piloting
		obu obu.OnboardUnit
		// the main autopilot loop and control structures
		work    func()
		done    chan bool
		trap    func(chan os.Signal)
		running bool
	}
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "autopilot")
}

// NewInstance creates and initializes a new autopilot instance
func NewInstance(obu obu.OnboardUnit) (*Autopilot, error) {
	ap := &Autopilot{
		obu:  obu,
		work: nil,
		done: make(chan bool, 1),
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
		running: false,
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

	return err
}

// Activate will start the autopilot
func (ap *Autopilot) Activate() error {
	logger.Info("Activating the autopilot")

	if ap.work == nil {
		ap.work = func() {}
	}
	go func() {
		ap.work()
		<-ap.done
	}()
	ap.running = true

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		ap.Stop()
	}()

	for {
		time.Sleep(1 * time.Second) // FIXME check if this is ok
		if !ap.running {
			break
		}
	}

	return nil
}

// Stop cancels the autopilot and stops the vehicle
func (ap *Autopilot) Stop() error {
	logger.Debug("stop")

	// FIXME do autopilot stuff here

	// signal the end of execution
	ap.done <- true
	ap.running = false

	return nil
}

// Shutdown finalizes the autopilot and releases all resources
func (ap *Autopilot) Shutdown() error {
	logger.Debug("shutdown")

	if ap.running {
		ap.Stop()
	}

	err := ap.obu.Shutdown()
	if err != nil {
		// FIXME something else?
		return err
	}
	// FIXME do autopilot stuff here
	return err
}

// AddWork is the main activity loop of the autopilot
func (ap *Autopilot) AddWork(f func()) error {
	ap.work = f
	return nil
}
