package telemetry

type (
	// DataFrame is a generic 'envelop' for sending telemetry data to the cloud
	DataFrame struct {
		DeviceID string
		Batch    int64
		Order    int64
		TS       int64
		Data     map[string]string
	}
)
