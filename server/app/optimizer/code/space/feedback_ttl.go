package space

import (
	"sync"
	"time"
)

const recentSize = 10

type TTL struct {
	recent []time.Duration
	mutex  sync.Mutex
}

func NewTTL() *TTL {
	return &TTL{recent: make([]time.Duration, 0), mutex: sync.Mutex{}}
}

func (t *TTL) Add(d time.Duration) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if len(t.recent) >= recentSize {
		recent := t.recent[:recentSize-1]
		t.recent = append([]time.Duration{d}, recent...)
		return
	}
	t.recent = append([]time.Duration{d}, t.recent...)
}

func (t *TTL) Time() time.Duration {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	tt := 0.0
	for i := 0; i < len(t.recent); i++ {
		tt += t.recent[i].Seconds()
	}

	if tt == 0.0 {
		return time.Duration(1e9)
	}
	tt /= float64(len(t.recent))

	return time.Duration(2 * 1e9 * tt)
}
