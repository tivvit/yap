package structs

import (
	"log"
	"strings"
)

type PipelineBlock interface {
	Run(state State)
	Checkable
}

type Checkable interface {
	Changed(state State, p *Pipeline) bool
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
			p.DepsFull = append(p.DepsFull, d)
		} else {
			if db, ok := p.Parent.Pipeline[d]; ok {
				switch db.(type) {
				case *Block:
					p.DepsFull = append(p.DepsFull, db.(*Block).FullName)
				case *Pipeline:
					p.DepsFull = append(p.DepsFull, db.(*Pipeline).FullName)
				}
			} else {
				log.Println("non-module dependency", d)
				p.DepsFull = append(p.DepsFull, d)
			}
		}
	}
}
