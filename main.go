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
	// todo user input
	log.Println("map", p.Map)
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
	// todo try to call dot
	p.Visualize()
	//p.Run(js)
	// todo generate report
	// todo generate state
	// todo generate
}
