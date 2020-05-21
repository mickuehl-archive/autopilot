package obu

import (
	"fmt"
	"shadow-racer/autopilot/v1/pkg/telemetry"
	"strconv"

	"github.com/majordomusio/commons/pkg/util"
	log "github.com/majordomusio/log15"
)

type (
	// OnboardUnit provides an interface to actual hardware
	OnboardUnit interface {
		// Initialize prepares the device
		Initialize() error
		// Reset re-initializes the device
		Reset() error
		// Shutdown releases all resources
		Shutdown() error
	}

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
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "obu")
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
		TS:       util.Timestamp(),
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
