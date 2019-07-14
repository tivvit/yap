package structs

import (
	"fmt"
	"log"
)

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

func stages(p PipelineBlock, prefix string) []string {
	var r []string
	switch p.(type) {
	case Pipeline:
		for k, v := range p.(Pipeline) {
			var name string
			if prefix != "" {
				name = fmt.Sprintf("%s_%s", prefix, k)
			} else {
				name = k
			}
			switch v.(type) {
			case Pipeline:
				r = append(r, stages(v, name)...)
			case *Pipeline:
				r = append(r, stages(*v.(*Pipeline), name)...)
			}
			r = append(r, name)
		}
	default:
		log.Fatal("Block has no stages")
	}
	return r
}

func (p Pipeline) Stages() []string {
	return stages(p, "")
}

func (p Pipeline) Plan() {
	log.Println(p.Stages())
	//g = graph.New()
}
