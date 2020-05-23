package parts

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/majordomusio/commons/pkg/util"

	"shadow-racer/autopilot/v1/pkg/eventbus"
	"shadow-racer/autopilot/v1/pkg/obu"
	"shadow-racer/autopilot/v1/pkg/telemetry"
)

const (
	stateDriving = "DRIVING"
	stateStopped = "STOPPED"
)

type (

	// Vehicle holds state information of a generic vehicle
	Vehicle struct {
		Mode        string  `json:"mode"` // some string, e.g. DRIVING, STOPPED, etc
		Throttle    float32 `json:"th"`   // -100 .. 100
		Steering    float32 `json:"st"`   // in deg, 0 is straight ahead
		Heading     float32 `json:"head"` // heading of the vehicle 0 -> North, 90 -> East ...
		Recording   bool    `json:"recording"`
		RecordingTS int64   `json:"recording_ts"`
		TS          int64   `json:"ts"` // timestamp
	}

	// VehicleState is an aggregate of vehicle state and other state information
	VehicleState struct {
		mutex   *sync.Mutex
		obu     obu.OnboardUnit
		vehicle *Vehicle
	}
)

// NewVehicleState creates a new state instance
func NewVehicleState(o obu.OnboardUnit) *VehicleState {
	return &VehicleState{
		mutex: &sync.Mutex{},
		obu:   o,
		vehicle: &Vehicle{
			Mode:        stateStopped,
			Throttle:    0.0,
			Steering:    0.0,
			Heading:     0.0,
			Recording:   false,
			RecordingTS: 0,
			TS:          util.TimestampNano(),
		},
	}
}

// Initialize prepares the device/component
func (v *VehicleState) Initialize() error {
	go v.RemoteStateHandler()
	return nil
}

// Reset re-initializes the device/component
func (v *VehicleState) Reset() error {
	return nil
}

// Shutdown releases all resources/component
func (v *VehicleState) Shutdown() error {
	return nil
}

// RemoteStateHandler listens remote state changes and updates the vehicle state accordingly
func (v *VehicleState) RemoteStateHandler() {
	logger.Info("Starting the remote state handler", "rxv", topicRCStateReceive, "txv", topicRCStateSend)

	ch := eventbus.InstanceOf().Subscribe(topicRCStateReceive)
	for {
		evt := <-ch
		state := evt.Data.(RemoteState)

		v.mutex.Lock()

		if state.Mode != v.vehicle.Mode {
			if state.Mode == stateDriving {
				// assumes v.vehicle.Mode == STOPPED
				//o.TailLights(4000, true) // FIXME enable tail lights
			} else if state.Mode == stateStopped {
				//o.TailLightsOff() // FIXME disable tail lights
				v.vehicle.Throttle = 0.0
				v.vehicle.Steering = 0.0
			} else {
				// FIXME should not happen
			}
			v.vehicle.Mode = state.Mode
			v.vehicle.Recording = state.Recording
		} else {
			//v.vehicle.Steering = 100.0 * ((float32(o.servo.MaxRange) / 90.0) * state.Steering)
			v.vehicle.Steering = 100.0 * ((30.0 / 90.0) * state.Steering) // FIXME -> o.servo.MaxRange config
			v.vehicle.Throttle = 100.0 * state.Throttle
		}

		if state.Recording != v.vehicle.Recording {
			baseURL := "http://localhost:3001" // FIXME configuration

			if state.Recording == true {
				v.vehicle.RecordingTS = util.Timestamp()
				v.vehicle.Recording = true

				resp, err := http.Get(fmt.Sprintf("%s/start?ts=%d", baseURL, v.vehicle.RecordingTS))

				if err != nil {
					logger.Error("Error toggling recording", "err", err.Error())
				} else {
					logger.Info("Started recording", "ts", v.vehicle.RecordingTS)
				}
				defer resp.Body.Close()
			} else {
				v.vehicle.Recording = false
				resp, err := http.Get(baseURL + "/stop")

				if err != nil {
					logger.Error("Error toggling recording", "err", err.Error())
				} else {
					logger.Info("Stopped recording")
				}
				defer resp.Body.Close()
			}
		}

		v.vehicle.TS = util.TimestampNano()

		// publish the new state
		eventbus.InstanceOf().Publish(topicRCStateSend, v.vehicle.Clone())

		// set the actuators
		v.obu.Direction(int(v.vehicle.Steering))
		v.obu.Throttle(int(v.vehicle.Throttle))

		v.mutex.Unlock()
	}
}

// NoopRemoteStateHandler is a no-op state handler i.e. values are just passed through
func NoopRemoteStateHandler() {
	logger.Info("Starting the NOOP remote state handler", "rxv", topicRCStateReceive, "txv", topicRCStateSend)

	ch := eventbus.InstanceOf().Subscribe(topicRCStateReceive)
	for {
		evt := <-ch
		state := evt.Data.(RemoteState)

		vehicle := Vehicle{
			Mode:        state.Mode,
			Steering:    100 * ((ServoRange / 90.0) * state.Steering),
			Throttle:    100 * state.Throttle,
			Heading:     360,
			Recording:   state.Recording,
			RecordingTS: util.TimestampNano(),
			TS:          util.TimestampNano(),
		}

		eventbus.InstanceOf().Publish(topicRCStateSend, &vehicle)
	}
}

// Clone returns a deep copy the vehicle state
func (v *Vehicle) Clone() *Vehicle {
	return &Vehicle{
		Mode:        v.Mode,
		Throttle:    v.Throttle,
		Steering:    v.Steering,
		Heading:     v.Heading,
		Recording:   v.Recording,
		RecordingTS: v.RecordingTS,
		TS:          v.TS,
	}
}

// ToDataFrame converts a vehicle state struct into a dataframe
func (v *Vehicle) ToDataFrame() *telemetry.DataFrame {
	df := telemetry.DataFrame{
		DeviceID: "shadow-racer",
		Batch:    v.RecordingTS,
		N:        v.TS,
		TS:       util.TimestampNano(),
		Type:     telemetry.KV,
		Data:     make(map[string]string),
	}
	df.Data["mode"] = v.Mode
	df.Data["th"] = fmt.Sprintf("%f", v.Throttle)
	df.Data["st"] = fmt.Sprintf("%f", v.Steering)
	df.Data["head"] = fmt.Sprintf("%f", v.Heading)
	df.Data["recording_ts"] = fmt.Sprintf("%d", v.RecordingTS)
	df.Data["ts"] = fmt.Sprintf("%d", v.TS)

	return &df
}

// ToVehicle creates an instance of vehicle state
func ToVehicle(df *telemetry.DataFrame) *Vehicle {

	if df.Type != telemetry.KV {
		return nil
	}

	v := Vehicle{
		Mode:      df.Data["mode"],
		Recording: true,
	}
	f, _ := strconv.ParseFloat(df.Data["th"], 32)
	v.Throttle = float32(f)
	f, _ = strconv.ParseFloat(df.Data["st"], 32)
	v.Steering = float32(f)
	f, _ = strconv.ParseFloat(df.Data["head"], 32)
	v.Heading = float32(f)
	i, _ := strconv.ParseInt(df.Data["ts"], 10, 64)
	v.TS = i
	i, _ = strconv.ParseInt(df.Data["recording_ts"], 10, 64)
	v.RecordingTS = i

	return &v
}
