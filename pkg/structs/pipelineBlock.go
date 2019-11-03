package structs

import (
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/stateStorage"
	log "github.com/sirupsen/logrus"
	"strings"
)

type PipelineBlock interface {
	Checkable
	Graphable
	Visualizable
	Run(state stateStorage.State, p *Pipeline)
	GetParent() *Pipeline
}

type Checkable interface {
	Changed(state stateStorage.State, p *Pipeline) bool
	GetState() (string, error)
}

type Graphable interface {
	GetDepsFull() []string
	GetFullName() string
}

type Visualizable interface {
	Visualize(ctx *dot.Graph, p *Pipeline, fileMap *map[string]*File, m *map[string]dot.Node, conf conf.VisualizeConf)
}

type PipelineBlockBase struct {
	Deps     []string  `yaml:"deps"`
	Parent   *Pipeline `yaml:"-"`
	FullName string    `yaml:"-"`
	DepsFull []string  `yaml:"-"`
}

func (p *PipelineBlockBase) genDepFull() {
	if len(p.DepsFull) > 0 {
		return
	}
	for _, d := range p.Deps {
		if strings.HasPrefix(d, "/") {
			// absolute dependency
			p.DepsFull = append(p.DepsFull, d)
		} else {
			// local (relative) dependency
			if db, ok := p.Parent.Pipeline[d]; ok {
				p.DepsFull = append(p.DepsFull, db.GetFullName())
			} else {
				log.Printf("Invalid dependency %s for %s\n", d, p.FullName)
			}
			//else {
			//	log.Println("non-module dependency", d)
			//
			//	p.DepsFull = append(p.DepsFull, d)
			//}
		}
	}
}

func (p PipelineBlockBase) GetFullName() string {
	return p.FullName
}

func (p PipelineBlockBase) GetParent() *Pipeline {
	return p.Parent
}
