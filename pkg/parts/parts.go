package parts

import (
	log "github.com/majordomusio/log15"

	"shadow-racer/autopilot/v1/pkg/pilot"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "parts")
}

/*

func (p *RaspiPilot) Direction(value int) {
	if value < p.cfg.servoRangeLimit*-1 {
		p.steering = p.cfg.servoRangeLimit * -1
	} else if value > p.cfg.servoRangeLimit {
		p.steering = p.cfg.servoRangeLimit
	} else {
		p.steering = value
	}

	// value == 0 -> servoMaxRange / 2
	direction := (p.cfg.servoMaxRange / 2) + p.steering + p.steeringTrim
	pulse := p.cfg.servoMinPulse + ((p.cfg.servoMaxPulse-p.cfg.servoMinPulse)/p.cfg.servoMaxRange)*direction

	// set the servo pulse
	//log.Printf("Direction: %d => (%d,%d)", p.steering, 0, pulse) // FIXME remove this
	err := p.actuators.SetPWM(p.cfg.servoChan, uint16(0), uint16(pulse))
	if err != nil {
		log.Printf(err.Error())
	}

}
*/

// StandardServoDirection sets the steering angle [-45,+45]
func StandardServoDirection(obu *pilot.OnboardUnit, value int) {
	logger.Debug("StandardServoDirection", "deg", value)
}

// StandardESCThrottle sets the motor speed [-100,+100]
func StandardESCThrottle(obu *pilot.OnboardUnit, value int) {
	logger.Debug("StandardESCThrottle", "thr", value)
}
