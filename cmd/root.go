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
	visualizeCmd.Flags().StringP("outImage", "i", "graph.png", "Output graph image filename")
	visualizeCmd.Flags().BoolP("noOutConn", "c", false, "Output file connections")
	visualizeCmd.Flags().BoolP("noPipelineNodes", "n", false, "Pipeline nodes")
	visualizeCmd.Flags().BoolP("noPipelineBoxes", "b", false, "PipelineBoxes")
	visualizeCmd.Flags().BoolP("noRunDot", "d", false, "Run dot")
	visualizeCmd.Flags().BoolP("noLegend", "l", false, "Display legend")
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
