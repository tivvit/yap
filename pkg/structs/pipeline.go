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
		for _, i := range ts {
			r = append(r, stageMap[nodeMap[i]])
		}
	}
	return r
}

func (p Pipeline) visualize(di *dot.Graph, pipeline PipelineBlock) (map[string]dot.Node, [][2]string) {
	nodesMap := map[string]dot.Node{}
	var deps [][2]string
	switch pipeline.(type) {
	case *Pipeline:
		for k, v := range pipeline.(*Pipeline).Pipeline {
			switch v.(type) {
			case *Block:
				//label := fmt.Sprintf("<table border=\"1\">%s<br/>%s</table>", k, strings.Join(v.(*Block).Exec, " "))
				label := fmt.Sprintf("<table BORDER=\"0\" CELLBORDER=\"1\" CELLSPACING=\"0\"><tr><td><b>%s</b></td></tr><tr><td><font face=\"Consolas\">%s</font></td></tr></table>", k, strings.Join(v.(*Block).Exec, " "))
				n := di.Node(k).Attr("shape", "plain")
				n.Attr("label", dot.HTML(label))
				nodesMap[k] = n
				for _, d := range v.(*Block).Deps {
					deps = append(deps, [2]string{d, k})
				}
				for _, o := range v.(*Block).Out {
					ob := di.Node(o)
					nodesMap[o] = ob
					log.Println(o, k)
					deps = append(deps, [2]string{o, k})
					log.Println(deps)
				}
			case *Pipeline:
				sg := di.Subgraph(k, dot.ClusterOption{})
				node := sg.Node(strings.ToUpper(k)).Attr("shape", "parallelogram")
				nodesMap[k] = node
				nodes, d := p.visualize(sg, v)
				for _, d := range v.(*Pipeline).Deps {
					deps = append(deps, [2]string{d, k})
				}
				for name, node := range nodes {
					nodesMap[name] = node
					// detect outputs
					att := node.AttributesMap
					if att.Value("shape") != nil {
						deps = append(deps, [2]string{name, k})
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
	nodeMap, deps := p.visualize(di, &p)
	for _, d := range deps {
		log.Println(d)
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
		log.Println("graphviz ok")
	}
}
