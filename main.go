package main

import (
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/stateStorage"
	"gopkg.in/yaml.v3"
	"log"
)

func main() {
	js := stateStorage.NewJsonStorage()
	p := pkg.Load()
	b, err := yaml.Marshal(p)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(b))
	// todo user input
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
