package pkg

import (
	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"github.com/tivvit/yap/cmdFlags"
	"github.com/tivvit/yap/pkg/structs"
	"github.com/tivvit/yap/pkg/utils"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
)

func potentialYapFile(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func findYapFile() string {
	var fn string
	fn = utils.MainFile
	if potentialYapFile(fn) {
		return fn
	}
	fn = utils.SecondaryMainFile
	if potentialYapFile(fn) {
		return fn
	}
	fn = utils.LastMainFile
	if potentialYapFile(fn) {
		return fn
	}
	return ""
}

func LoadCmd(cmd *cobra.Command) *structs.Pipeline {
	f, err := cmd.Flags().GetString(cmdFlags.File)
	if err != nil {
		log.Fatalln(err)
	}
	return Load(f)
}

func Load(fileName string) *structs.Pipeline {
	var yapFile string
	if fileName == "" {
		yapFile = findYapFile()
	} else {
		yapFile = fileName
	}
	b, err := ioutil.ReadFile(yapFile)
	if err != nil {
		log.Printf("load %s err #%v ", yapFile, err)
	}
	p := structs.PipelineRaw{}
	//log.Println(string(b))
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	pp := ParsePipeline(&p)
	pp.Enrich()
	return pp
}

func LoadFile(fileName string) *structs.PipelineRaw {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("load %s err #%v ", fileName, err)
	}
	p := structs.PipelineRaw{}
	//log.Println(string(b))
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &p
}

func LoadScript(name string) (*structs.PipelineRaw, error) {
	c, err := shellwords.Parse(name)
	if err != nil {
		log.Fatalln(err)
	}
	out, ok := utils.GenericRunEnv(c, []string{}, false, false)
	if !ok {
		log.Fatalf("`%s` failed during loading pipeline", c)
	}
	p := structs.PipelineRaw{}
	err = yaml.Unmarshal([]byte(out), &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &p, nil
}
