package event

import "time"

var (
	BlockRunEventName = "block-run-event"
)

type BlockRunEvent struct {
	BaseEvent `yaml:",inline"`
	Block     string         `json:"block,omitempty"`
	Duration  *time.Duration `json:"duration-nano,omitempty"`
	StartTime *time.Time     `json:"start-time,omitempty"`
}

func NewBlockRunEvent(message string, block string) *BlockRunEvent {
	e := &BlockRunEvent{
		BaseEvent: *NewEvent(message),
	}
	e.Type = BlockRunEventName
	e.Block = block
	e.AddTag("block")
	e.AddTag("run")
	return e
}
