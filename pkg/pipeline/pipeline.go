package pipeline

import (
	"github.com/tivvit/yap/pkg/structs"
	"github.com/yourbasic/graph"
	"log"
)

func Filter(name string, stageMap map[string]structs.Graphable) map[string]structs.Graphable {
	if name != "" {
		stageMap = filterDeps(stageMap, name)
	} else {
		stageMap = filterMain(stageMap)
	}
	return stageMap
}

func filterDeps(stageMap map[string]structs.Graphable, name string) map[string]structs.Graphable {
	var found bool
	_, found = stageMap[name]
	if !found {
		// try files
		oldName := name
		name = structs.DotFilePrefix + name
		log.Printf("%s not found trying %s", oldName, name)
		_, found = stageMap[name]
	}
	if !found {
		log.Printf("Target %s not found \n", name)
		return map[string]structs.Graphable{}
	}
	g, m, mi := CreateInverseGraph(stageMap)
	r := make(map[string]structs.Graphable)
	r[name] = stageMap[name]
	graph.BFS(g, m[name], func(f, t int, _ int64) {
		log.Printf("%s -> %s (%d -> %d)", mi[f], mi[t], f, t)
		r[mi[f]] = stageMap[mi[f]]
		r[mi[t]] = stageMap[mi[t]]
	})
	return r
}

func filterMain(sm map[string]structs.Graphable) map[string]structs.Graphable {
	fsm := map[string]structs.Graphable{}
	for k, s := range sm {
		switch s.(type) {
		case structs.PipelineBlock:
			if s.(structs.PipelineBlock).GetParent().GetFullName() == structs.MainFullName {
				fsm[k] = s
			}
		}
	}
	return fsm
}

func CreateInverseGraph(stages map[string]structs.Graphable) (*graph.Mutable, map[string]int, map[int]string) {
	return createGraph(stages, true)
}

func CreateGraph(stages map[string]structs.Graphable) (*graph.Mutable, map[string]int, map[int]string) {
	return createGraph(stages, false)
}

func createGraph(stages map[string]structs.Graphable, inverse bool) (*graph.Mutable, map[string]int, map[int]string) {
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
			nn, ok := m[n]
			if !ok {
				log.Printf("Elem %s not present in map", n)
			}
			dd, ok := m[d]
			if !ok {
				log.Printf("Elem %s not present in map", d)
			}
			if inverse {
				g.Add(nn, dd)
			} else {
				g.Add(dd, nn)
			}
		}
	}
	return g, m, mi
}

func Plan(p *structs.Pipeline, name string) []structs.PipelineBlock {
	stageMap := Filter(name, p.GetGraphable())
	g, _, mi := CreateGraph(stageMap)
	if !graph.Acyclic(g) {
		log.Println("There is a cycle in the dependencies")
		return []structs.PipelineBlock{}
	}
	ts, _ := graph.TopSort(g)
	var r []structs.PipelineBlock
	for _, v := range ts {
		e := stageMap[mi[v]]
		switch e.(type) {
		case structs.PipelineBlock:
			r = append(r, e.(structs.PipelineBlock))
		}
	}
	return r
}