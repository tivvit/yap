package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// todo support taint (force run)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("RUN")
		//for _, s := range pl {
		//	s.Run(js, p)
		//}

		// todo generate report
		// todo generate state
		// todo generate
	},
}
