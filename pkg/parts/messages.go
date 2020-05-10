package parts

type (
	// RemoteState is the current driving state info
	RemoteState struct {
		Mode string `json:"mode"`

		Throttle float32 `json:"th"`
		Steering float32 `json:"st"`

		TS int64 `json:"ts"`
	}

	// HUD holds display info sent to a remote client
	HUD struct {
		Mode string `json:"mode,omitempty"`

		Throttle float32 `json:"th"`   // 0 .. 100
		Steering float32 `json:"st"`   // in deg
		Heading  float32 `json:"head"` // heading of the vehicle 0 -> North, 90 -> East ...

		TS int64 `json:"ts"` // timestamp
	}
)
