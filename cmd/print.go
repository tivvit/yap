package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/pkg"
	"gopkg.in/yaml.v3"
	log "github.com/sirupsen/logrus"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Aliases: []string{"p"},
	Short: "print final yaml",
	Long:  ``,
	Version: yapCmd.Version,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		p := pkg.LoadCmd(cmd)
		b, err := yaml.Marshal(p)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(string(b))
	},
}