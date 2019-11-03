package event

import (
	"strings"
	"time"
)

var (
	BaseEventName = "generic-event"
)

type Event interface{}


type BaseEvent struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Message   string    `json:"message,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	Type      string    `json:"type"`
}

func NewEvent(message string) *BaseEvent {
	return &BaseEvent{
		Timestamp: time.Now(),
		Message:   message,
		Type:      BaseEventName,
	}
}

func (e *BaseEvent) AddTag(tag string) {
	e.Tags = append(e.Tags, strings.ToLower(tag))
}
