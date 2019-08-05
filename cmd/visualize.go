package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/stateStorage"
	"gopkg.in/yaml.v3"
	"log"
)

var visualizeCmd = &cobra.Command{
	Use:   "vis",
	Short: "visualize",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		o, err := cmd.Flags().GetString("out")
		log.Println("OUTPUT FILE", o)

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
		js := stateStorage.NewJsonStorage()
		//log.Println("parent", p.Parent)
		//log.Println(p.Pipeline["test"].(*structs.Pipeline).DepsFull)
		//log.Println(p.Pipeline["finalize"].(*structs.Block).DepsFull)
		//pl := p.Plan("/finalize")
		//log.Println(pl)
		//log.Println(p.Plan("/main/B"))
		//log.Println(p.Plan("/main/a"))
		//log.Println(p.Plan("/main/A"))
		//log.Println(p.Plan("files.txt"))
		//log.Println(p.Plan(""))
		//b, err = yaml.Marshal(pl)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//log.Println(string(b))
		log.Println(js)
		// todo support directory dependency
		// todo check missing deps
		//p.Vis("/finalize")

		if len(args) == 0 {
			p.Vis("")
		} else if len(args) == 1 {
			p.Vis(args[0])
		} else {
			log.Fatalln("too many args")
		}
	},
}
