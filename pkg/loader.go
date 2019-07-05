package pkg

import (
	"bytes"
	"github.com/tivvit/yap/pkg/structs"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os/exec"
)

func Load() *structs.Pipeline {
	b, err := ioutil.ReadFile(MainFile)
	if err != nil {
		log.Printf("load yap.yml err #%v ", err)
	}
	p := structs.PipelineRaw{}
	//log.Println(string(b))
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return ParsePipeline(&p)
}

func LoadFile(fileName string) *structs.PipelineRaw {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("load yap.yml err #%v ", err)
	}
	p := structs.PipelineRaw{}
	//log.Println(string(b))
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return &p
}

func LoadScript(name string) *structs.PipelineRaw {
	cmd := exec.Command("python3", name)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%v", out.String())
	p := structs.PipelineRaw{}
	err = yaml.Unmarshal(out.Bytes(), &p)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return &p
}