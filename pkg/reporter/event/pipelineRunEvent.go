package event

import "time"

var (
	PipelineRunEventName = "pipeline-run-event"
)

type PipelineRunEvent struct {
	BaseEvent `yaml:",inline"`
	Pipeline  string         `json:"pipeline,omitempty"`
	Duration  *time.Duration `json:"duration-nano,omitempty"`
	StartTime *time.Time     `json:"start-time,omitempty"`
}

func NewPipelineRunEvent(message string, pipeline string) *PipelineRunEvent {
	e := &PipelineRunEvent{
		BaseEvent: *NewEvent(message),
	}
	e.Type = PipelineRunEventName
	e.Pipeline = pipeline
	e.AddTag("pipeline")
	e.AddTag("run")
	return e
}
