package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/structs"
	"log"
)

var visualizeCmd = &cobra.Command{
	Use:   "vis",
	Short: "visualize",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		o, err := cmd.Flags().GetString("out")
		if err != nil {
			log.Fatalln(err)
		}
		oi, err := cmd.Flags().GetString("outImage")
		if err != nil {
			log.Fatalln(err)
		}
		oc, err := cmd.Flags().GetBool("noOutConn")
		if err != nil {
			log.Fatalln(err)
		}
		pn, err := cmd.Flags().GetBool("noPipelineNodes")
		if err != nil {
			log.Fatalln(err)
		}
		pb, err := cmd.Flags().GetBool("noPipelineBoxes")
		if err != nil {
			log.Fatalln(err)
		}
		d, err := cmd.Flags().GetBool("noRunDot")
		if err != nil {
			log.Fatalln(err)
		}
		l, err := cmd.Flags().GetBool("noLegend")
		if err != nil {
			log.Fatalln(err)
		}

		conf := structs.VisualizeConf{
			OutputFile: o,
			OutputImage: oi,
			OutputConnections: !oc,
			PipelineNodes: !pn,
			PipelineBoxes: !pb,
			RunDot: !d,
			Legend: !l,
		}

		log.Println(conf)

		p := pkg.Load()
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
			pkg.Visualize(p, "", conf)
		} else if len(args) == 1 {
			pkg.Visualize(p, args[0], conf)
		} else {
			log.Fatalln("too many args")
		}
	},
}
