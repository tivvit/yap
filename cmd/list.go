package cmd

import (
	"fmt"
	"regexp"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
)

var listCmd = &cobra.Command{
	Use:     "list [name-regex]",
	Aliases: []string{"l"},
	Short:   "list",
	Long:    ``,
	Version: yapCmd.Version,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		p := pkg.LoadCmd(cmd)

		if len(args) == 0 {
			printList(p.List(), func(s string) bool { return true })
		} else if len(args) == 1 {
			userRe := args[0]
			re, err := regexp.Compile(userRe)
			if err != nil {
				log.Fatalf("Filter regex \"%s\" is invalid\n", userRe)
			}
			printList(p.List(), re.MatchString)
		} else {
			log.Fatalln("too many args")
		}
	},
}

func printList(b []string, f func(s string) bool) {
	// todo add formatted output option
	sort.Strings(b)
	for _, i := range b {
		if f(i) {
			fmt.Println(i)
		}
	}
}
