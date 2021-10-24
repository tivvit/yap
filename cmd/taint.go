package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
	"regexp"
)

var taintCmd = &cobra.Command{
	Use:     "taint name-regex",
	Aliases: []string{"t"},
	Short:   "taint block",
	Long:    ``,
	Version: yapCmd.Version,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userRe := args[0]
		re, err := regexp.Compile(userRe)
		if err != nil {
			log.Fatalf("Filter regex \"%s\" is invalid\n", userRe)
		}

		p := pkg.LoadCmd(cmd)
		js := p.State
		for _, i := range p.List() {
			if re.MatchString(i) {
				js.Delete(i)
				log.Infof("Tainting %s", i)
			}
		}
	},
}
