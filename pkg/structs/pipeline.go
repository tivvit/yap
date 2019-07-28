package structs

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/yourbasic/graph"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Pipeline struct {
	*PipelineBlockBase `yaml:",inline"`
	Version            float32                  `yaml:"version"`
	Settings           map[string]interface{}   `yaml:"settings,omitempty"`
	Pipeline           map[string]PipelineBlock `yaml:"pipeline"`
	Map                map[string]PipelineBlock `yaml:"-"`
	MapFiles           map[string]*File         `yaml:"-"`
}

func NewPipeline(version float32, settings map[string]interface{}, deps []string) *Pipeline {
	return &Pipeline{
		Version:  version,
		Settings: settings,
		Pipeline: make(map[string]PipelineBlock),
		MapFiles: make(map[string]*File),
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

func (p Pipeline) Run(state stateStorage.State, pl *Pipeline) {
	for k, v := range p.Pipeline {
		log.Printf("Running %s", k)
		v.Run(state, pl)
	}
}

func (pl Pipeline) Changed(state stateStorage.State, p *Pipeline) bool {
	for _, b := range p.Pipeline {
		if b.Changed(state, p) {
			return true
		}
	}
	return false
}

func (p Pipeline) GetState() (string, error) {
	m := make(map[string]string)
	for _, b := range p.Pipeline {
		st, err := b.GetState()
		if err != nil {
			m[b.GetFullName()] = ""
		} else {
			m[b.GetFullName()] = st
		}
	}
	jb, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	ch := string(jb)
	//ch, err := utils.Md5Checksum(bytes.NewReader(jb))
	//if err != nil {
	//	return "", err
	//}
	return ch, nil
}

func (p Pipeline) GetDepsFull() []string {
	return p.DepsFull
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

func (p Pipeline) addFile(name string) *File {
	if f, ok := p.MapFiles[name]; ok {
		return f
	} else {
		p.MapFiles[name] = &File{
			Name: name,
		}
		p.MapFiles[name].Analyze()
		return p.MapFiles[name]
	}
}

func (p *Pipeline) files() {
	var f *File
	for _, d := range p.Map {
		switch d.(type) {
		case *Block:
			for _, o := range d.(*Block).Out {
				f = p.addFile(o)
				f.Deps = append(f.Deps, d.(*Block))
			}
			for _, i := range d.(*Block).In {
				f = p.addFile(i)
			}
		}
	}
}

func (p *Pipeline) Enrich() {
	p.names("", p)
	p.parents(p)
	p.Map = p.flatten(p)
	p.genDepFullRec(p)
	p.files()
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
				// todo this can cause name clash = it should not due to namespacing
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

func (p Pipeline) GetGraphable() map[string]Graphable {
	m := make(map[string]Graphable)
	for k, v := range p.Map {
		m[k] = v
	}
	for k, v := range p.MapFiles {
		if _, ok := m[k]; ok {
			log.Printf("%s is duplicate (present in block and in files) - will use the block", k)
		} else {
			m[k] = v
		}
	}
	return m
}

func CreateInverseGraph(stages map[string]Graphable) (*graph.Mutable, map[string]int) {
	return createGraph(stages, true)
}

func CreateGraph(stages map[string]Graphable) (*graph.Mutable, map[string]int) {
	return createGraph(stages, false)
}

func createGraph(stages map[string]Graphable, inverse bool) (*graph.Mutable, map[string]int) {
	g := graph.New(len(stages))
	m := make(map[string]int)
	i := 0
	for k := range stages {
		m[k] = i
		i++
	}
	for n, s := range stages {
		for _, d := range s.GetDepsFull() {
			if inverse {
				g.Add(m[n], m[d])
			} else {
				g.Add(m[d], m[n])
			}
		}
	}
	return g, m
}

func (p Pipeline) Plan(name string) []PipelineBlock {
	stageMap := p.GetGraphable()
	g, m := CreateInverseGraph(stageMap)
	mi := make(map[int]string)
	for k, v := range m {
		mi[v] = k
	}
	target, found := m[name]
	if !found {
		log.Printf("Target %s not found \n", name)
		return []PipelineBlock{}
	}
	graph.BFS(g, target, func(f, t int, _ int64) {
		log.Printf("%s -> %s (%d -> %d)", mi[f], mi[t], f, t)
	})
	log.Println("---")
	return []PipelineBlock{}

	//nodeMap := make(map[int]string)
	//nodeMapInv := make(map[string]int)
	//c := 0
	//for k := range stageMap {
	//	nodeMap[c] = k
	//	nodeMapInv[k] = c
	//	c++
	//}
	//log.Println(stageMap)
	//log.Println(p.GetGraphable())
	//log.Println(nodeMap)
	//log.Println(nodeMapInv)
	//g := graph.New(len(stageMap))
	//gi := graph.New(len(stageMap))
	//for k, v := range stageMap {
	//	switch v.(type) {
	//	case *Pipeline:
	//		for _, d := range v.(*Pipeline).Deps {
	//			if id, ok := nodeMapInv[d]; ok {
	//				g.Add(id, nodeMapInv[k])
	//				gi.Add(nodeMapInv[k], id)
	//				log.Println("adding edge", d, "->", k, id, "->", nodeMapInv[k])
	//			} else {
	//				log.Println(d, "is not a target")
	//			}
	//		}
	//	case *Block:
	//		for _, d := range v.(*Block).Deps {
	//			if id, ok := nodeMapInv[d]; ok {
	//				g.Add(id, nodeMapInv[k])
	//				gi.Add(nodeMapInv[k], id)
	//				log.Println("adding edge", d, "->", k, id, "->", nodeMapInv[k])
	//			} else {
	//				log.Println(d, "is not a target")
	//			}
	//		}
	//	}
	//
	//}
	//// todo extract only parents - from bfs of inverted graph
	//log.Println("Acyclic", graph.Acyclic(g))
	//log.Println(g.String())
	//ts, ok := graph.TopSort(g)
	//if !ok {
	//	log.Println("Graph contains cycle")
	//	return nil
	//}
	//log.Println(ts)
	//var r []PipelineBlock
	//if name != "" {
	//	r = append(r, stageMap[name])
	//	graph.BFS(gi, nodeMapInv[name], func(f, t int, _ int64) {
	//		//log.Println(f, "to", t)
	//		r = append(r, stageMap[nodeMap[t]])
	//	})
	//	for left, right := 0, len(r)-1; left < right; left, right = left+1, right-1 {
	//		r[left], r[right] = r[right], r[left]
	//	}
	//} else {
	//	for _, i := range ts {
	//		r = append(r, stageMap[nodeMap[i]])
	//	}
	//}
	//return r
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
				nodesMap[b.FullName] = n
			case *Pipeline:
				sg := di.Subgraph(k, dot.ClusterOption{})
				node := sg.Node(strings.ToUpper(k)).Attr("shape", "parallelogram")
				nodesMap[v.(*Pipeline).FullName] = node
				nodes, d := p.visualize(sg, main, v)
				for name, node := range nodes {
					nodesMap[name] = node
					//log.Println(name)
					// detect outputs
					att := node.AttributesMap
					// todo file deps for pipelines are unsupported = delete this
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
	fileMap := make(map[string]dot.Node)
	for _, f := range p.MapFiles {
		if f.Analyzed && f.IsDir {
			fileMap[f.Name] = di.Node(f.Name).Attr("shape", "septagon")
		} else {
			fileMap[f.Name] = di.Node(f.Name).Attr("shape", "oval")
		}
	}
	for _, n := range p.Map {
		switch n.(type) {
		case *Block:
			b := n.(*Block)
			for _, t := range b.DepsFull {
				log.Println(t)
				di.Edge(nodeMap[t], nodeMap[b.FullName])
			}
			for _, o := range b.Out {
				di.Edge(nodeMap[b.FullName], fileMap[o])
			}
			for _, i := range b.In {
				di.Edge(fileMap[i], nodeMap[b.FullName])
			}
		case *Pipeline:
			b := n.(*Pipeline)
			log.Println(b.DepsFull)
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
