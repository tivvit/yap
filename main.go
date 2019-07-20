package main

import (
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/structs"
	"gopkg.in/yaml.v3"
	"log"
)

// todo support taint (force run)
// todo cli options

func main() {
	js := stateStorage.NewJsonStorage()
	p := pkg.Load()
	b, err := yaml.Marshal(p)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(b))
	//for k := range p.Map {
	//	log.Println(k)
	//}
	// todo user input
	log.Println("map", p.Map)
	log.Println("mapFiles")
	for k, f := range p.MapFiles {
		log.Println(k, f.Name, f.Deps)
		for _, d := range f.Deps {
			log.Println(d.FullName)
		}
	}
	log.Println("parent", p.Parent)
	log.Println(p.Pipeline["test"].(*structs.Pipeline).DepsFull)
	log.Println(p.Pipeline["finalize"].(*structs.Block).DepsFull)
	pl := p.Plan("finalize")
	//pl := p.Plan("")
	b, err = yaml.Marshal(pl)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(b))
	log.Println(js)
	// todo support directory dependency
	p.Visualize()
	p.Run(js, p)
	// todo generate report
	// todo generate state
	// todo generate
}
