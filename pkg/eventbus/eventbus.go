package eventbus

// see https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997

import (
	"sync"

	log "github.com/majordomusio/log15"
)

type (
	// DataEvent wraps a message being published to subscribers
	DataEvent struct {
		Data  interface{}
		Topic string
	}

	// DataChannel is a channel which can accept an DataEvent
	DataChannel chan DataEvent

	// DataChannelSlice is a slice of DataChannels
	DataChannelSlice []DataChannel

	// EventBus stores the information about subscribers interested for a particular topic
	EventBus struct {
		subscribers map[string]DataChannelSlice
		rm          sync.RWMutex
	}
)

var (
	logger   log.Logger
	eventbus *EventBus
)

func init() {
	logger = log.New("module", "eventbus")

	eb := &EventBus{
		subscribers: map[string]DataChannelSlice{},
	}
	eventbus = eb
}

// InstanceOf enforces a singleton pattern
func InstanceOf() *EventBus {
	return eventbus
}

// Publish sends any kind of data to a topic
func (eb *EventBus) Publish(topic string, data interface{}) {
	eb.rm.RLock()

	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		// special thanks for /u/freesid who pointed it out
		channels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
	eb.rm.RUnlock()
}

// Subscribe returns a channel that listens on topic
func (eb *EventBus) Subscribe(topic string) DataChannel {
	eb.rm.Lock()

	ch := make(chan DataEvent)
	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
	eb.rm.Unlock()

	return ch
}
