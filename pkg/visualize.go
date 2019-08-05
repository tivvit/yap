package pkg

import (
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/pipeline"
	"github.com/tivvit/yap/pkg/structs"
	"log"
	"os"
	"os/exec"
)

func Visualize(p *structs.Pipeline, name string, conf structs.VisualizeConf) {
	sm := pipeline.Filter(name, p.GetGraphable())
	di := dot.NewGraph(dot.Directed)
	nodeMap := map[string]dot.Node{}
	fileMap := map[string]*structs.File{}
	for _, v := range sm {
		switch v.(type) {
		case structs.Visualizable:
			v.(structs.Visualizable).Visualize(di, &fileMap, &nodeMap, conf)
		default:
			log.Println("Unexpected type %T which is not Visualizable", v)
		}
	}
	for n := range fileMap {
		fileMap[n] = p.MapFiles[n]
		// todo duplicate names with blocks
		fileMap[n].Visualize(di, &fileMap, &nodeMap, conf)
	}
	for k, n := range p.Map {
		if _, ok := nodeMap[k]; !ok {
			continue
		}
		switch n.(type) {
		case structs.PipelineBlock:
			for _, d := range n.(structs.PipelineBlock).GetDepsFull() {
				di.Edge(nodeMap[d], nodeMap[n.GetFullName()])
			}
		}
		if conf.OutputConnections {
			switch n.(type) {
			case *structs.Block:
				for _, d := range n.(*structs.Block).Out {
					di.Edge(nodeMap[n.GetFullName()], nodeMap[d])
				}
			}
		}
		if conf.PipelineNodes {
			switch n.(type) {
			case *structs.Pipeline:
				for _, d := range n.(*structs.Pipeline).Pipeline {
					di.Edge(nodeMap[d.GetFullName()], nodeMap[n.GetFullName()])
				}
			}
		}
	}

	if conf.Legend {
		legend(di, conf)
	}

	f, _ := os.Create(conf.OutputFile)
	di.Write(f)

	if conf.RunDot {
		tryDot(conf)
	}
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

func legend(di *dot.Graph, conf structs.VisualizeConf) {
	legend := di.Subgraph("Legend", dot.ClusterOption{})
	legend.Node("File").Attr("shape", structs.FileShape)
	legend.Node("Directory").Attr("shape", structs.DirShape)
	pipeline := legend
	if conf.PipelineBoxes {
		pipeline = legend.Subgraph("Pipeline", dot.ClusterOption{})
	}
	if conf.PipelineNodes {
		pipeline.Node("Pipeline").Attr("shape", structs.PipelineShape)
	}
	b := structs.Block{
		Name:        "Block",
		Description: "This is a Block",
		Exec:        []string{"code", "-h"},
	}
	legendMap := map[string]dot.Node{}
	b.Visualize(legend, nil, &legendMap, conf)
}

func tryDot(conf structs.VisualizeConf) {
	c := exec.Command("dot", []string{"-T", "png", conf.OutputFile, "-o", conf.OutputImage}...)
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Graphviz ok")
	}
}
