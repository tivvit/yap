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
	p.Plan()
	log.Println(js)
	//p.Run(js)
	// todo resolve deps
	// todo generate report
	// todo generate state
	// todo generate
}
