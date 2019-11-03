package tracker

import (
	"errors"
	"fmt"
	"time"
)

type tracker struct {
	streams map[string]time.Time
}

func NewTracker() *tracker {
	return &tracker{
		streams: map[string]time.Time{},
	}
}

func (t tracker) Start(name string) {
	t.streams[name] = time.Now()
}

func (t* tracker) Stop(name string) (time.Duration, time.Time, error) {
	if v, ok := t.streams[name]; ok {
		delete(t.streams, name)
		return time.Since(v), v, nil
	} else {
		return time.Duration(0), time.Time{}, errors.New(fmt.Sprintf("tracker %s not started", name))
	}
}