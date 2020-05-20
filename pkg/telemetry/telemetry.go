package telemetry

type (
	// DataMap is a generic name/value pair map with all strings
	DataMap map[string]string

	// DataFrame is a generic 'envelop' for sending telemetry data to the cloud
	DataFrame struct {
		DeviceID string  `json:"deviceid"`
		Batch    int64   `json:"batch,omitempty"`
		Order    int64   `json:"order,omitempty"`
		TS       int64   `json:"ts"`
		Data     DataMap `json:"data,omitempty"`
	}
)
