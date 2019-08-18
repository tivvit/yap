package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"github.com/tivvit/yap/pkg"
	"github.com/tivvit/yap/pkg/stateStorage"
	"log"
)


var runCmd = &cobra.Command{
	Use:   "run",
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
		// todo support filter
		p.Run(js, p)

		// todo support taint (force run)
		// todo generate report
		// todo generate state
		// todo generate
	},
}
