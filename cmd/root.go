package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
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


	},
}

func init() {
	yapCmd.AddCommand(runCmd)
	yapCmd.AddCommand(visualizeCmd)
	yapCmd.AddCommand(printCmd)
	yapCmd.PersistentFlags().StringP(cmdFlags.File, "f", "", "Main yapfile path")
	runCmd.Flags().BoolP(cmdFlags.DryRun, "d", false, "Do not run - just check")
	visualizeCmd.Flags().StringP(cmdFlags.Out, "o", "graph.dot", "Output graph filename")
	visualizeCmd.Flags().StringP(cmdFlags.OutImage, "i", "graph.png", "Output graph image filename")
	visualizeCmd.Flags().BoolP(cmdFlags.NoOutConn, "C", false, "Disable output file connections")
	visualizeCmd.Flags().BoolP(cmdFlags.NoPipelineNodes, "N", false, "Disable pipeline nodes")
	visualizeCmd.Flags().BoolP(cmdFlags.NoPipelineBoxes, "B", false, "Disable pipeline boxes")
	visualizeCmd.Flags().BoolP(cmdFlags.NoRunDot, "D", false, "Do not Run dot")
	visualizeCmd.Flags().BoolP(cmdFlags.NoLegend, "L", false, "Do not Display legend")
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
