package structs

import (
	"github.com/tivvit/yap/pkg/stateStorage"
	"log"
	"strings"
)

type PipelineBlock interface {
	Run(state stateStorage.State, p *Pipeline)
	Checkable
	GetFullName() string
}

type Checkable interface {
	Changed(state stateStorage.State, p *Pipeline) bool
	GetState() (string, error)
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
		// todo this may be file with absolute path
		if strings.HasPrefix(d, "/") {
			p.DepsFull = append(p.DepsFull, d)
		} else {
			if db, ok := p.Parent.Pipeline[d]; ok {
				p.DepsFull = append(p.DepsFull, db.GetFullName())
			} else {
				log.Println("non-module dependency", d)
				p.DepsFull = append(p.DepsFull, d)
			}
		}
	}
}

func (p PipelineBlockBase) GetFullName() string {
	return p.FullName
}
