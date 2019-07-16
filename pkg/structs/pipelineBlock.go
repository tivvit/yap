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
	Changed(state State) bool
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
	// todo deps for files
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
				log.Println("non module dependency", d)
				// todo this is probably file
			}
		}
	}
}
