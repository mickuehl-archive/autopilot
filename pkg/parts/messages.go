package parts

type (
	// RemoteState is the current driving state info
	RemoteState struct {
		Throttle float32 `json:"th,omitempty"`
		Steering float32 `json:"st,omitempty"`
		Mode     string  `json:"mode,omitempty"`
		TS       int64   `json:"ts,omitempty"`
	}

	// HUD holds display info sent to a remote client
	HUD struct {
		Throttle float32 // 0 .. 100
		Steering float32 // in deg
		Heading  float32 // heading of the vehicle 0 -> North, 90 -> East ...
		TS       int64   // timestamp
	}
)
