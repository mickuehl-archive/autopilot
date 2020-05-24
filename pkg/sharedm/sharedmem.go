package sharedm

import (
	"sync"
)

type (
	memBytes struct {
		data []byte
		mux  *sync.RWMutex
	}

	sharedMemory struct {
		mux         *sync.RWMutex
		sharedBytes map[string]memBytes
	}
)

var (
	memory sharedMemory
)

func init() {
	memory = sharedMemory{
		mux:         &sync.RWMutex{},
		sharedBytes: make(map[string]memBytes),
	}
}

// StoreInt stores an integer value v with key k
func StoreInt(k string, v int64) {

}

// StoreFloat stores a float value v with key k
func StoreFloat(k string, v float64) {

}

// StoreString stores a pointer to string v with key k
func StoreString(k, v *string) {

}

// StoreBytes stores a pointer to bytes v with key k
func StoreBytes(k string, v []byte) {
	if b, ok := memory.sharedBytes[k]; ok {
		b.mux.Lock()
		b.data = v
		b.mux.Unlock()
	} else {
		memory.mux.Lock()
		b := memBytes{
			data: v,
			mux:  &sync.RWMutex{},
		}
		memory.sharedBytes[k] = b
		memory.mux.Unlock()
	}
}

// GetInt returns value of key k. The bool returned indicates hit/miss.
func GetInt(k string) (int64, bool) {
	return 0, false
}

// GetFloat returns value of key k. The bool returned indicates hit/miss.
func GetFloat(k string) (float64, bool) {
	return 0, false
}

// GetString returns value of key k. The bool returned indicates hit/miss.
func GetString(k string) (*string, bool) {
	s := ""
	return &s, false
}

// GetBytes returns value of key k. The bool returned indicates hit/miss.
func GetBytes(k string) ([]byte, bool) {
	if b, ok := memory.sharedBytes[k]; ok {
		b.mux.RLock()
		defer b.mux.RUnlock()
		return b.data, true
	} else {
		return nil, false
	}
}
