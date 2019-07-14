package structs

import (
	"fmt"
	"github.com/yourbasic/graph"
	"log"
)

type PipelineBlock interface {
	Run(state State)
}

type Pipeline struct {
	Version  float32                  `yaml:"version"`
	Settings map[string]interface{}   `yaml:"settings,omitempty"`
	Pipeline map[string]PipelineBlock `yaml:"pipeline"`
	Deps     []string                 `yaml:"deps"`
}

func NewPipeline(version float32, settings map[string]interface{}, deps []string) *Pipeline {
	return &Pipeline{
		Version:  version,
		Settings: settings,
		Deps:     deps,
		Pipeline: make(map[string]PipelineBlock),
	}
}

type PipelineRaw struct {
	Version  float32                `yaml:"version"`
	Settings map[string]interface{} `yaml:"settings,omitempty"`
	Pipeline map[string]interface{} `yaml:"pipeline"`
	Deps     []string               `yaml:"deps"`
}

func (p Pipeline) Run(state State) {
	for k, v := range p.Pipeline {
		log.Printf("Running %s", k)
		v.Run(state)
	}
}

func stages(p PipelineBlock, prefix string) map[string]PipelineBlock {
	r := make(map[string]PipelineBlock)
	switch p.(type) {
	case *Pipeline:
		for k, v := range p.(*Pipeline).Pipeline {
			var name string
			if prefix != "" {
				name = fmt.Sprintf("%s_%s", prefix, k)
			} else {
				name = k
			}
			switch v.(type) {
			case *Pipeline:
				// todo this can cause name clash
				for n, block := range stages(v.(*Pipeline), name) {
					r[n] = block
				}
				r[k] = v
			case *Block:
				r[name] = v
			default:
				log.Println("Unexpected type")
			}
		}
	default:
		log.Fatal("Block has no stages")
	}
	return r
}

func (p Pipeline) Stages() map[string]PipelineBlock {
	return stages(&p, "")
}

func (p Pipeline) Plan(name string) []PipelineBlock {
	// todo plan per block (not everything together)
	stageMap := p.Stages()
	nodeMap := make(map[int]string)
	nodeMapInv := make(map[string]int)
	c := 0
	for k := range stageMap {
		nodeMap[c] = k
		nodeMapInv[k] = c
		c++
	}
	log.Println(stageMap)
	log.Println(nodeMap)
	log.Println(nodeMapInv)
	g := graph.New(len(stageMap))
	gi := graph.New(len(stageMap))
	for k, v := range stageMap {
		switch v.(type) {
		case *Pipeline:
			for _, d := range v.(*Pipeline).Deps {
				if id, ok := nodeMapInv[d]; ok {
					g.Add(id, nodeMapInv[k])
					gi.Add(nodeMapInv[k], id)
					log.Println("adding edge", d, "->", k, id, "->", nodeMapInv[k])
				} else {
					log.Println(d, "is not a target")
				}
			}
		case *Block:
			for _, d := range v.(*Block).Deps {
				if id, ok := nodeMapInv[d]; ok {
					g.Add(id, nodeMapInv[k])
					gi.Add(nodeMapInv[k], id)
					log.Println("adding edge", d, "->", k, id, "->", nodeMapInv[k])
				} else {
					log.Println(d, "is not a target")
				}
			}
		}

	}
	log.Println("Acyclic", graph.Acyclic(g))
	log.Println(g.String())
	ts, _ := graph.TopSort(g)
	log.Println(ts)
	var r []PipelineBlock
	if name != "" {
		r = append(r, stageMap[name])
		graph.BFS(gi, nodeMapInv[name], func(f, t int, _ int64) {
			//log.Println(f, "to", t)
			r = append(r, stageMap[nodeMap[t]])
		})
		for left, right := 0, len(r)-1; left < right; left, right = left+1, right-1 {
			r[left], r[right] = r[right], r[left]
		}
	} else {
		for _, i  := range ts {
			r = append(r, stageMap[nodeMap[i]])
		}
	}
	return r
}
