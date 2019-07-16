package structs

import (
	"fmt"
	"github.com/emicklei/dot"
	"github.com/yourbasic/graph"
	"log"
	"os"
	"os/exec"
	"strings"
)

// todo files may be output and input
// todo dependencies for files

type Pipeline struct {
	*PipelineBlockBase `yaml:",inline"`
	Version            float32                  `yaml:"version"`
	Settings           map[string]interface{}   `yaml:"settings,omitempty"`
	Pipeline           map[string]PipelineBlock `yaml:"pipeline"`
	Map                map[string]PipelineBlock `yaml:"-"`
}

func NewPipeline(version float32, settings map[string]interface{}, deps []string) *Pipeline {
	return &Pipeline{
		Version:  version,
		Settings: settings,
		Pipeline: make(map[string]PipelineBlock),
		PipelineBlockBase: &PipelineBlockBase{
			Deps: deps,
		},
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

func (p Pipeline) Changed(state State) bool {
	for _, b := range p.Pipeline {
		if b.Changed(state) {
			return true
		}
	}
	return false
}

func (p *Pipeline) genDepFullRec(pb PipelineBlock) {
	switch pb.(type) {
	case *Block:
		pb.(*Block).genDepFull()
	case *Pipeline:
		for _, v := range pb.(*Pipeline).Pipeline {
			p.genDepFullRec(v)
		}
		pb.(*Pipeline).genDepFull()
	}
}

func getFullName(name string, namespace string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func (p Pipeline) names(namespace string, pb PipelineBlock) {
	switch pb.(type) {
	case *Pipeline:
		for k, v := range pb.(*Pipeline).Pipeline {
			switch v.(type) {
			case *Pipeline:
				v.(*Pipeline).FullName = getFullName(k, namespace)
				p.names(getFullName(k, namespace), v.(*Pipeline))
			case *Block:
				v.(*Block).FullName = getFullName(k, namespace)
			}
		}
	}
}

func (p Pipeline) parents(pb PipelineBlock) {
	switch pb.(type) {
	case *Pipeline:
		for _, v := range pb.(*Pipeline).Pipeline {
			switch v.(type) {
			case *Pipeline:
				p.parents(v.(*Pipeline))
				v.(*Pipeline).Parent = pb.(*Pipeline)
			case *Block:
				v.(*Block).Parent = pb.(*Pipeline)
			}
		}
	}
}

func (p Pipeline) flatten(pb PipelineBlock) map[string]PipelineBlock {
	r := map[string]PipelineBlock{}
	switch pb.(type) {
	case *Pipeline:
		for _, v := range pb.(*Pipeline).Pipeline {
			switch v.(type) {
			case *Pipeline:
				for n, b := range p.flatten(v.(*Pipeline)) {
					r[n] = b
				}
				r[v.(*Pipeline).FullName] = v.(*Pipeline)
			case *Block:
				r[v.(*Block).FullName] = v.(*Block)
			}
		}
	}
	return r
}

func (p *Pipeline) Enrich() {
	p.names("", p)
	p.parents(p)
	p.Map = p.flatten(p)
	p.genDepFullRec(p)
}

// todo outputs with blocks references
// todo search dep method

func stages(p PipelineBlock, prefix string) map[string]PipelineBlock {
	r := make(map[string]PipelineBlock)
	switch p.(type) {
	case *Pipeline:
		for k, v := range p.(*Pipeline).Pipeline {
			var name string
			if prefix != "" {
				name = fmt.Sprintf("%s/%s", prefix, k)
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
		for _, i := range ts {
			r = append(r, stageMap[nodeMap[i]])
		}
	}
	return r
}

func (p Pipeline) visualize(di *dot.Graph, main *dot.Graph, pipeline PipelineBlock) (map[string]dot.Node, [][2]string) {
	nodesMap := map[string]dot.Node{}
	var deps [][2]string
	switch pipeline.(type) {
	case *Pipeline:
		for k, v := range pipeline.(*Pipeline).Pipeline {
			switch v.(type) {
			case *Block:
				b := v.(*Block)
				nameFmt := "<tr><td><b>%s</b></td></tr>"
				name := fmt.Sprintf(nameFmt, k)
				cmdFmt := `<tr><td><font face="Courier New, Courier, monospace">%s</font></td></tr>`
				cmd := fmt.Sprintf(cmdFmt, strings.Join(b.Exec, " "))
				descFmt := "<tr><td>%s</td></tr>"
				desc := ""
				if b.Description != "" {
					desc = fmt.Sprintf(descFmt, b.Description)
				}
				tableFmt := `<table border="0" cellborder="1" cellspacing="0">%s%s%s</table>`
				label := fmt.Sprintf(fmt.Sprintf(tableFmt, name, desc, cmd))
				n := di.Node(k).Attr("shape", "plain")
				n.Attr("label", dot.HTML(label))
				nodesMap[v.(*Block).FullName] = n
				//for _, d := range v.(*Block).Deps {
				//	deps = append(deps, [2]string{d, k})
				//}
				// todo add file deps
				for _, o := range v.(*Block).Out {
					ob := main.Node(o)
					nodesMap[o] = ob
					deps = append(deps, [2]string{v.(*Block).FullName, o})
				}
			case *Pipeline:
				sg := di.Subgraph(k, dot.ClusterOption{})
				node := sg.Node(strings.ToUpper(k)).Attr("shape", "parallelogram")
				nodesMap[v.(*Pipeline).FullName] = node
				nodes, d := p.visualize(sg, main, v)
				//for _, d := range v.(*Pipeline).Deps {
				//	deps = append(deps, [2]string{d, k})
				//}
				for name, node := range nodes {
					nodesMap[name] = node
					log.Println(name)
					// detect outputs
					att := node.AttributesMap
					if att.Value("shape") != nil {
						deps = append(deps, [2]string{name, v.(*Pipeline).FullName})
					}
				}
				for _, i := range d {
					deps = append(deps, i)
				}
			}
		}
	}
	return nodesMap, deps
}

func (p Pipeline) Visualize() {
	di := dot.NewGraph(dot.Directed)
	// todo file separation should be optional
	nodeMap, deps := p.visualize(di, di, &p)
	for _, n := range p.Map {
		switch n.(type) {
		case *Block:
			b := n.(*Block)
			for _, t := range b.DepsFull {
				di.Edge(nodeMap[t], nodeMap[b.FullName])
			}
		case *Pipeline:
			b := n.(*Pipeline)
			for _, t := range b.DepsFull {
				di.Edge(nodeMap[t], nodeMap[b.FullName])
			}
		}

	}
	for _, d := range deps {
		di.Edge(nodeMap[d[0]], nodeMap[d[1]])
	}
	f, _ := os.Create("graph.dot")
	di.Write(f)
	p.tryDot()
}

func (p Pipeline) tryDot() {
	c := exec.Command("dot", []string{"-T", "png", "graph.dot", "-o", "graph.png"}...)
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Graphviz ok")
	}
}
