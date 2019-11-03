package event

import "time"

var (
	RunEventName = "run-event"
)

type RunEvent struct {
	BaseEvent `yaml:",inline"`
	Command   string         `json:"command,omitempty"`
	Env       string         `json:"env,omitempty"`
	Duration  *time.Duration `json:"duration-nano,omitempty"`
	StartTime *time.Time     `json:"start-time,omitempty"`
}

func NewRunEvent(message string) *RunEvent {
	e := &RunEvent{
		BaseEvent: *NewEvent(message),
	}
	e.Type = RunEventName
	e.AddTag("run")
	return e
}
