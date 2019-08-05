package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// todo support taint (force run)

var yapCmd = &cobra.Command{
	Use:   "yap",
	Short: "Yet Another Pipeline",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			log.Fatalln(err)
		}
		os.Exit(0)

		//for _, s := range pl {
		//	s.Run(js, p)
		//}

		// todo generate report
		// todo generate state
		// todo generate
	},
}

func init() {
	yapCmd.AddCommand(runCmd)
	yapCmd.AddCommand(visualizeCmd)
	visualizeCmd.Flags().StringP("out", "o", "graph.dot", "Output graph filename")
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
