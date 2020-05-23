package telemetry

const (
	BLOB = iota
	KV
)

type (
	// DataMap is a generic name/value pair map with all strings
	DataMap map[string]string

	// DataFrame is a generic 'envelop' for sending telemetry data to the cloud
	DataFrame struct {
		DeviceID string  `json:"deviceid"`        // any identifier of the device, e.g. a hostname, serial number etc.
		Batch    int64   `json:"batch,omitempty"` // a prefix to group telemetry data
		Type     int     `json:"type"`            // the type of the payload: BLOB or KV
		Data     DataMap `json:"data,omitempty"`  // data in a key/value set
		Blob     string  `json:"blob,omitempty"`  // binary data in base64 encoded string
		TS       int64   `json:"ts"`              // the timestamp the message was sent
	}
)
