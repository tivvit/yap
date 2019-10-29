package event

import "time"

type Event struct {
	Timestamp time.Time      `json:"timestamp,omitempty"`
	Block     string         `json:"block,omitempty"`
	Duration  *time.Duration `json:"duration,omitempty"`
	StartTime *time.Time     `json:"start-time,omitempty"`
	EndTime   *time.Time     `json:"end-time,omitempty"`
	Message   string         `json:"message,omitempty"`
	Tags      []string       `json:"tags,omitempty"`
}

func NewEvent() *Event {
	return &Event{
		Timestamp: time.Now(),
	}
}
