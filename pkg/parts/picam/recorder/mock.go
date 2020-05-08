/*
 * Copyright (C) 2018 Josh A. Beam
 * All rights reserved.
 *
 * See https://github.com/joshb/pi-camera-go
 */

package recorder

import (
	"time"
)

type mockRecorder struct {
	running     bool
	subscribers []Subscriber
}

func NewMock() Recorder {
	return &mockRecorder{}
}

func (r *mockRecorder) Start() error {
	r.running = true
	return nil
}

func (r *mockRecorder) Stop() error {
	r.running = false
	return nil
}

func (r *mockRecorder) SegmentDuration() time.Duration {
	return 5 * time.Second
}

func (r *mockRecorder) AddSubscriber(subscriber Subscriber) {
	r.subscribers = append(r.subscribers, subscriber)
}
