package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/pipeline"
	"github.com/tivvit/yap/pkg/reporter"
	"github.com/tivvit/yap/pkg/stateStorage"
	log "github.com/sirupsen/logrus"
)


var runCmd = &cobra.Command{
	Use:   "run [block-name]",
	Aliases: []string{"r"},
	Short: "run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		p := pkg.LoadCmd(cmd)
		js := stateStorage.NewJsonStorage()
		dr, err := cmd.Flags().GetBool(cmdFlags.DryRun)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(dr)
		// todo dry run

		reporter.NewReporter(conf.ReporterConf{
			Storages: []conf.ReporterStorageConf{
				conf.ReporterStorageConfJson{
					FileName: "report.json",
				},
			},
		})

		if len(args) == 0 {
			p.Run(js, p)
		} else if len(args) == 1 {
			pl := pipeline.Plan(p, args[0])
			log.Println(pl)
			for _, ps := range pl {
				ps.Run(js, p)
			}
		} else {
			log.Fatalln("too many args")
		}

		// todo support taint (force run)
		// todo generate report
		// todo generate state
		// todo generate
	},
}
