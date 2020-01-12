package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"os"
)

// todo support taint (force run)

var (
	GitCommit string
	GitTag    = "unknown"
)

var yapCmd = &cobra.Command{
	Use:     "yap",
	Short:   "Yet Another Pipeline",
	Long:    ``,
	Args:    cobra.OnlyValidArgs,
	Version: GitTag + " " + GitCommit,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setLogger(cmd)
	},
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
	yapCmd.AddCommand(listCmd)
	yapCmd.AddCommand(taintCmd)
	yapCmd.PersistentFlags().StringP(cmdFlags.File, "f", "", "Main yapfile path")
	yapCmd.PersistentFlags().BoolP(cmdFlags.Quiet, "q", false, "Suppress yap output")
	//yapCmd.PersistentFlags().BoolP(cmdFlags.Quiet, "q", false, "Suppress yap output")
	runCmd.Flags().BoolP(cmdFlags.DryRun, "d", false, "Do not run - just check")
	visualizeCmd.Flags().StringP(cmdFlags.Out, "o", "graph.dot", "Output graph filename")
	visualizeCmd.Flags().StringP(cmdFlags.OutImage, "i", "graph.png", "Output graph image filename")
	visualizeCmd.Flags().BoolP(cmdFlags.NoOutConn, "C", false, "Disable output file connections")
	visualizeCmd.Flags().BoolP(cmdFlags.NoPipelineNodes, "N", false, "Disable pipeline nodes")
	visualizeCmd.Flags().BoolP(cmdFlags.NoPipelineBoxes, "B", false, "Disable pipeline boxes")
	visualizeCmd.Flags().BoolP(cmdFlags.NoRunDot, "D", false, "Do not Run dot")
	visualizeCmd.Flags().BoolP(cmdFlags.NoLegend, "L", false, "Do not Display legend")
	visualizeCmd.Flags().BoolP(cmdFlags.Check, "s", false, "Check state changes")
}

func setLogger(cmd *cobra.Command) {
	quiet, err := cmd.Flags().GetBool(cmdFlags.Quiet)
	if err != nil {
		log.Warnln("Quiet flag not loaded")
	}
	if quiet {
		log.SetLevel(log.FatalLevel)
	}
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
