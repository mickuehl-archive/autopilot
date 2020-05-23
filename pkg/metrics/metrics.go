package metrics

import (
	log "github.com/majordomusio/log15"
)

var (
	logger  log.Logger
	metrics map[string]Meter
)

func init() {
	logger = log.New("module", "metrics")
	metrics = make(map[string]Meter)
}

// NewMeter constructs a new StandardMeter and launches a goroutine
func NewMeter(item string) Meter {

	m := newStandardMeter()

	arbiter.Lock()
	defer arbiter.Unlock()
	arbiter.meters[m] = struct{}{}
	if !arbiter.started {
		arbiter.started = true
		go arbiter.tick()
	}

	// add to the registry
	if m1, ok := metrics[item]; ok {
		m1.Stop()
	}
	metrics[item] = m

	return m
}

// DumpMeters dumps all meters to the log
func DumpMeters() {
	for k, m := range metrics {
		logger.Debug("Meter", "meter", k, "c", m.Count(), "mean", m.RateMean(), "r1", m.Rate1(), "r5", m.Rate5(), "r15", m.Rate15())
	}
}

// Mark measures the number & frequency of something
func Mark(item string) {
	if m, ok := metrics[item]; ok {
		m.Mark(1)
	}
}

// MarkN measures the number & frequency of something
func MarkN(item string, n int64) {
	if m, ok := metrics[item]; ok {
		m.Mark(n)
	}
}
