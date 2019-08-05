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
	Short: "",
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

func init() {
	yapCmd.AddCommand(runCmd)
	yapCmd.AddCommand(visualizeCmd)
	//visualizeCmd.Flags().StringP("source", "s", "", "Source directory to read from")
}

func Execute() {
	if err := yapCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
