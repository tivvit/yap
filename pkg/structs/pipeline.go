package structs

import "log"

type PipelineBlock interface {
	Run(state State)
}

type Pipeline struct {
	Version  float32                  `yaml:"version"`
	Settings map[string]interface{}   `yaml:"settings,omitempty"`
	Pipeline map[string]PipelineBlock `yaml:"pipeline"`
}

func NewPipeline(version float32, settings map[string]interface{}) *Pipeline {
	return &Pipeline{
		Version:  version,
		Settings: settings,
		Pipeline: make(map[string]PipelineBlock),
	}
}

type PipelineRaw struct {
	Version  float32                `yaml:"version"`
	Settings map[string]interface{} `yaml:"settings,omitempty"`
	Pipeline map[string]interface{} `yaml:"pipeline"`
}

func (p Pipeline) Run(state State) {
	for k, v := range p.Pipeline {
		log.Printf("Running %s", k)
		v.Run(state)
	}
}
