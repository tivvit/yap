package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
	"log"
	"sort"
)

var listCmd = &cobra.Command{
	Use:   "list [name-regex]",
	Short: "list",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		p := pkg.LoadCmd(cmd)

		if len(args) == 0 {
			printList(p.List())
		} else if len(args) == 1 {
			// todo regex filter
			printList(p.List())
		} else {
			log.Fatalln("too many args")
		}
	},
}

func printList(b []string) {
	sort.Strings(b)
	for _, i := range b {
		fmt.Println(i)
	}
}
