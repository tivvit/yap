package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"github.com/tivvit/yap/pkg"
	conf2 "github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/stateStorage"
)

var visualizeCmd = &cobra.Command{
	Use:     "visualize [block-name]",
	Aliases: []string{"v", "vis"},
	Short:   "visualize",
	Long:    ``,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		o, err := cmd.Flags().GetString(cmdFlags.Out)
		if err != nil {
			log.Fatalln(err)
		}
		oi, err := cmd.Flags().GetString(cmdFlags.OutImage)
		if err != nil {
			log.Fatalln(err)
		}
		oc, err := cmd.Flags().GetBool(cmdFlags.NoOutConn)
		if err != nil {
			log.Fatalln(err)
		}
		pn, err := cmd.Flags().GetBool(cmdFlags.NoPipelineNodes)
		if err != nil {
			log.Fatalln(err)
		}
		pb, err := cmd.Flags().GetBool(cmdFlags.NoPipelineBoxes)
		if err != nil {
			log.Fatalln(err)
		}
		d, err := cmd.Flags().GetBool(cmdFlags.NoRunDot)
		if err != nil {
			log.Fatalln(err)
		}
		l, err := cmd.Flags().GetBool(cmdFlags.NoLegend)
		if err != nil {
			log.Fatalln(err)
		}
		s, err := cmd.Flags().GetBool(cmdFlags.Check)
		if err != nil {
			log.Fatalln(err)
		}

		conf := conf2.VisualizeConf{
			OutputFile:        o,
			OutputImage:       oi,
			OutputConnections: !oc,
			PipelineNodes:     !pn,
			PipelineBoxes:     !pb,
			RunDot:            !d,
			Legend:            !l,
			Check:             s,
		}

		log.Println(conf)

		p := pkg.LoadCmd(cmd)
		//for k := range p.Map {
		//	log.Println(k)
		//}
		// todo user input

		//log.Println("map", p.Map)
		//log.Println("mapFiles", p.MapFiles)
		//for k, f := range p.MapFiles {
		//	log.Println(k, f.Name, f.Deps)
		//	for _, d := range f.Deps {
		//		log.Println(d.FullName)
		//	}
		//}

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
