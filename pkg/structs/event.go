package structs

import "time"

type Event struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	StartTime *time.Time `json:"start-time,omitempty"`
	EndTime   *time.Time `json:"end-time,omitempty"`
	Message   string     `json:"message,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
}

func NewEvent() *Event {
	return &Event{
		Timestamp: time.Now(),
	}
}
