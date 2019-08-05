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

const (
	mainNamespace = ""
	mainName      = ""
	mainFullName  = mainNamespace + "/" + mainName
	pipelineShape = "parallelogram"
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
	p.FullName = getFullName(mainName, namespace)
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
			lm := map[string]PipelineBlock{}
			switch v.(type) {
			case *Pipeline:
				for n, b := range p.flatten(v.(*Pipeline)) {
					r[n] = b
					lm[n] = b
				}
				v.(*Pipeline).Map = lm
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
	p.names(mainNamespace, p)
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

func CreateInverseGraph(stages map[string]Graphable) (*graph.Mutable, map[string]int, map[int]string) {
	return createGraph(stages, true)
}

func CreateGraph(stages map[string]Graphable) (*graph.Mutable, map[string]int, map[int]string) {
	return createGraph(stages, false)
}

func createGraph(stages map[string]Graphable, inverse bool) (*graph.Mutable, map[string]int, map[int]string) {
	g := graph.New(len(stages))
	m := make(map[string]int)
	mi := make(map[int]string)
	i := 0
	for k := range stages {
		m[k] = i
		mi[i] = k
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
	return g, m, mi
}

func filterDeps(stageMap map[string]Graphable, name string) map[string]Graphable {
	_, found := stageMap[name]
	if !found {
		log.Printf("Target %s not found \n", name)
		return map[string]Graphable{}
	}
	g, m, mi := CreateInverseGraph(stageMap)
	r := make(map[string]Graphable)
	graph.BFS(g, m[name], func(f, t int, _ int64) {
		log.Printf("%s -> %s (%d -> %d)", mi[f], mi[t], f, t)
		r[mi[f]] = stageMap[mi[f]]
		r[mi[t]] = stageMap[mi[t]]
	})
	return r
}

func filterMain(sm map[string]Graphable) map[string]Graphable {
	fsm := map[string]Graphable{}
	for k, s := range sm {
		switch s.(type) {
		case PipelineBlock:
			if s.(PipelineBlock).GetParent().GetFullName() == mainFullName {
				fsm[k] = s
			}
		}
	}
	return fsm
}

func (p Pipeline) Plan(name string) []PipelineBlock {
	stageMap := filter(name, p.GetGraphable())
	g, _, mi := CreateGraph(stageMap)
	if !graph.Acyclic(g) {
		log.Println("There is a cycle in the dependencies")
		return []PipelineBlock{}
	}
	ts, _ := graph.TopSort(g)
	var r []PipelineBlock
	for _, v := range ts {
		e := stageMap[mi[v]]
		switch e.(type) {
		case PipelineBlock:
			r = append(r, e.(PipelineBlock))
		}
	}
	return r
}

func filter(name string, stageMap map[string]Graphable) map[string]Graphable {
	if name != "" {
		stageMap = filterDeps(stageMap, name)
	} else {
		stageMap = filterMain(stageMap)
	}
	return stageMap
}

//func (p Pipeline) visualize(di *dot.Graph, main *dot.Graph, pipeline PipelineBlock) (map[string]dot.Node, [][2]string) {
//	nodesMap := map[string]dot.Node{}
//	var deps [][2]string
//	switch pipeline.(type) {
//	case *Pipeline:
//		for k, v := range pipeline.(*Pipeline).Pipeline {
//			switch v.(type) {
//			case *Block:
//				b := v.(*Block)
//				nameFmt := "<tr><td><b>%s</b></td></tr>"
//				name := fmt.Sprintf(nameFmt, k)
//				cmdFmt := `<tr><td><font face="Courier New, Courier, monospace">%s</font></td></tr>`
//				cmd := fmt.Sprintf(cmdFmt, strings.Join(b.Exec, " "))
//				descFmt := "<tr><td>%s</td></tr>"
//				desc := ""
//				if b.Description != "" {
//					desc = fmt.Sprintf(descFmt, b.Description)
//				}
//				tableFmt := `<table border="0" cellborder="1" cellspacing="0">%s%s%s</table>`
//				label := fmt.Sprintf(fmt.Sprintf(tableFmt, name, desc, cmd))
//				n := di.Node(k).Attr("shape", "plain")
//				n.Attr("label", dot.HTML(label))
//				nodesMap[b.FullName] = n
//			case *Pipeline:
//				sg := di.Subgraph(k, dot.ClusterOption{})
//				node := sg.Node(strings.ToUpper(k)).Attr("shape", "parallelogram")
//				nodesMap[v.(*Pipeline).FullName] = node
//				nodes, d := p.visualize(sg, main, v)
//				for name, node := range nodes {
//					nodesMap[name] = node
//					//log.Println(name)
//					// detect outputs
//					att := node.AttributesMap
//					// todo file deps for pipelines are unsupported = delete this
//					if att.Value("shape") != nil {
//						deps = append(deps, [2]string{name, v.(*Pipeline).FullName})
//					}
//				}
//				for _, i := range d {
//					deps = append(deps, i)
//				}
//			}
//		}
//	}
//	return nodesMap, deps
//}

func (p Pipeline) Visualize(ctx *dot.Graph, fileMap *map[string]*File, m *map[string]dot.Node) {
	splitName := strings.Split(p.FullName, "/")
	name := splitName[len(splitName)-1]
	log.Println("name is", name, "for", p.FullName)
	if name != "" {
		innerCtx := ctx
		//innerCtx := ctx.Subgraph(name, dot.ClusterOption{})

		// virtual block node
		node := innerCtx.Node(strings.ToUpper(name)).Attr("shape", pipelineShape)
		(*m)[p.FullName] = node

		//nodes, d := p.visualize(sg, main, v)
		//for name, node := range nodes {
		//	(*m)[name] = node
		//}
		for _, v := range p.Pipeline {
			v.Visualize(innerCtx, fileMap, m)
		}
	} else {
		for _, v := range p.Pipeline {
			v.Visualize(ctx, fileMap, m)
		}
	}
}

func (p Pipeline) Vis(name string) {
	sm := filter(name, p.GetGraphable())
	di := dot.NewGraph(dot.Directed)
	nodeMap := map[string]dot.Node{}
	fileMap := map[string]*File{}
	for _, v := range sm {
		switch v.(type) {
		case Visualizable:
			v.(Visualizable).Visualize(di, &fileMap, &nodeMap)
		default:
			log.Println("Unexpected type %T which is not Visualizable", v)
		}
	}
	for n := range fileMap {
		fileMap[n] = p.MapFiles[n]
		fileMap[n].Visualize(di, &fileMap, &nodeMap)
	}
	for k, n := range p.Map {
		if _, ok := nodeMap[k]; !ok {
			continue
		}
		switch n.(type) {
		case PipelineBlock:
			for _, d := range n.(PipelineBlock).GetDepsFull() {
				di.Edge(nodeMap[d], nodeMap[n.GetFullName()])
			}
		}
		switch n.(type) {
		case *Block:
			for _, d := range n.(*Block).Out {
				di.Edge(nodeMap[n.GetFullName()], nodeMap[d])
			}
		}
		// todo connect to parent block
		//switch n.(type) {
		//case *Pipeline:
		//	for _, d := range n.(*Pipeline).Pipeline {
		//		di.Edge(nodeMap[d.GetFullName()], nodeMap[n.GetFullName()])
		//	}
		//}
	}

	legend := di.Subgraph("Legend", dot.ClusterOption{})
	legend.Node("File").Attr("shape", fileShape)
	legend.Node("Directory").Attr("shape", dirShape)
	pipeline := legend.Subgraph("Pipeline", dot.ClusterOption{})
	pipeline.Node("Pipeline").Attr("shape", pipelineShape)
	b := Block{
		Name: "Block",
		Description: "This is a Block",
		Exec: []string{"code", "-h"},
	}
	legendMap := map[string]dot.Node{}
	b.Visualize(legend, nil, &legendMap)

	f, _ := os.Create("graph.dot")
	di.Write(f)
	p.tryDot()
	return

	//// todo file separation should be optional
	//nodeMap, deps := p.visualize(di, di, &p)
	//fileMap := make(map[string]dot.Node)
	//for _, f := range p.MapFiles {
	//	f.Visualize(di, &fileMap)
	//	//if f.Analyzed && f.IsDir {
	//	//	fileMap[f.Name] = di.Node(f.Name).Attr("shape", "septagon")
	//	//} else {
	//	//	fileMap[f.Name] = di.Node(f.Name).Attr("shape", "oval")
	//	//}
	//}
	//for _, n := range p.Map {
	//	switch n.(type) {
	//	case *Block:
	//		b := n.(*Block)
	//		for _, t := range b.DepsFull {
	//			log.Println(t)
	//			di.Edge(nodeMap[t], nodeMap[b.FullName])
	//		}
	//		for _, o := range b.Out {
	//			di.Edge(nodeMap[b.FullName], fileMap[o])
	//		}
	//		for _, i := range b.In {
	//			di.Edge(fileMap[i], nodeMap[b.FullName])
	//		}
	//	case *Pipeline:
	//		b := n.(*Pipeline)
	//		log.Println(b.DepsFull)
	//		for _, t := range b.DepsFull {
	//			di.Edge(nodeMap[t], nodeMap[b.FullName])
	//		}
	//	}
	//
	//}
	//for _, d := range deps {
	//	di.Edge(nodeMap[d[0]], nodeMap[d[1]])
	//}
	//f, _ := os.Create("graph.dot")
	//di.Write(f)
	//p.tryDot()
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
