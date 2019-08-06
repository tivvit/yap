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
	yapCmd.AddCommand(printCmd)
	visualizeCmd.Flags().StringP("out", "o", "graph.dot", "Output graph filename")
	visualizeCmd.Flags().StringP("out-image", "i", "graph.png", "Output graph image filename")
	visualizeCmd.Flags().BoolP("no-out-conn", "C", false, "Disable output file connections")
	visualizeCmd.Flags().BoolP("no-pipeline-nodes", "N", false, "Disable pipeline nodes")
	visualizeCmd.Flags().BoolP("no-pipeline-boxes", "B", false, "Disable pipeline boxes")
	visualizeCmd.Flags().BoolP("no-run-dot", "D", false, "Do not Run dot")
	visualizeCmd.Flags().BoolP("no-legend", "L", false, "Do not Display legend")
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
