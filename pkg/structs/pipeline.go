package structs

import "log"

type PipelineBlock interface {
	Run(state State)
}

type Pipeline map[string]PipelineBlock
type PipelineRaw map[string]interface{}

func (p Pipeline) Run(state State) {
	for k, v := range p {
		log.Printf("Running %s", k)
		v.Run(state)
	}
}
