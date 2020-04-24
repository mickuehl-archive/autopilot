package pilot

type (
	// Config holds basic pre-sets
	Config struct {
		frequency int // PCA9685 clock speed
		// actuators
		servoChan int // channel of the servo
		escChan   int // channel of the speed controller
		led1Chan  int // channels of the two back/break LEDs
		led2Chan  int

		// actuator pre-sets

		// min pulse length out of 4096
		servoMinPulse int
		// max pulse length out of 4096
		servoMaxPulse int
		// the max the servo can rotate (in deg)
		servoMaxRange int
		// the max the servo is allowed to rotate to each side (deg)
		servoRangeLimit int

		// min pulse length out of 4096, rel to escZeroPulse
		escMinPulse int
		// max pulse length out of 4096, rel to escZeroPulse
		escMaxPulse int
		// zero pulse length out of 4096
		escZeroPulse int
		// base pulse to work from
		escBasePulse int
		// is the pulse to reset/initialize the ESC
		escInitPulse int
		// esc can reverse if true
		escCanReverse bool
	}

	Pilot interface {
		// Start initializes the pilot and all its components
		Start() error
		// Shutdown stops all components
		Shutdown() error
		// Direction sets the steering angle
		Direction(value int)
		// Throttle sets the motor speed
		Throttle(value float32)
		// Turn the back lights on/off, w/o blinking
		BackLights(value, f int)
	}
)
