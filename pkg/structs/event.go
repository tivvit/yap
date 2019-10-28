package structs

import "time"

type Event struct {
	StartTime *time.Time `json:"start-time,omitempty"`
	EndTime   *time.Time `json:"end-time,omitempty"`
	Message   string     `json:"message,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
}
