package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/pipeline"
	"github.com/tivvit/yap/pkg/reporter"
)

var runCmd = &cobra.Command{
	Use:     "run [block-name]",
	Aliases: []string{"r"},
	Short:   "run",
	Long:    ``,
	Version: yapCmd.Version,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		p := pkg.LoadCmd(cmd)
		js := p.State
		dry, err := cmd.Flags().GetBool(cmdFlags.DryRun)
		if err != nil {
			log.Fatalln(err)
		}
		if dry {
			// todo configure reporter by flags (pipeline conf)
			log.Infoln("Dry run")
			reporter.NewReporter(conf.ReporterConf{
				Storages: []conf.ReporterStorageConf{},
			})
		} else {
			reporter.NewReporter(conf.ReporterConf{
				Storages: []conf.ReporterStorageConf{
					conf.ReporterStorageConfJson{
						FileName: "report.json",
					},
				},
			})
		}

		if len(args) == 0 {
			p.Run(js, p, dry)
		} else if len(args) == 1 {
			pl := pipeline.Plan(p, args[0])
			for _, ps := range pl {
				ps.Run(js, p, dry)
			}
		} else {
			log.Fatalln("too many args")
		}

		// todo generate report
		// todo generate state
		// todo generate
	},
}
