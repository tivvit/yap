package structs

import (
	"encoding/json"
	"errors"
	"github.com/mattn/go-shellwords"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"log"
	"strings"
)

type Block struct {
	PipelineBlockBase `yaml:",inline"`
	Name              string   `yaml:"-"`
	Description       string   `yaml:"desc,omitempty"`
	Check             []string `yaml:"check,omitempty"`
	Exec              []string `yaml:"exec,omitempty"`
	Out               []string `yaml:"out,omitempty"`
}

type IncludeBlock struct {
	Include string
	Type    string
}

func NewBlockFromMap(name string, m map[string]interface{}) *Block {
	b := Block{
		PipelineBlockBase: PipelineBlockBase{},
	}
	b.Name = name
	// todo keys names based on struct annotation
	if m["exec"] == nil {
		log.Fatal("exec field missing")
	}
	ex, err := shellwords.Parse(m["exec"].(string))
	if err != nil {
		log.Fatalln(err)
	}
	b.Exec = ex
	//strings.Split(m["exec"].(string), " ")
	if m["desc"] != nil {
		b.Description = m["desc"].(string)
	}
	if m["check"] != nil {
		// todo force string
		ex, err = shellwords.Parse(m["check"].(string))
		if err != nil {
			log.Fatalln(err)
		}
		b.Check = ex
	}
	if m["deps"] != nil {
		var deps []string
		for _, d := range m["deps"].([]interface{}) {
			deps = append(deps, d.(string))
		}
		b.Deps = deps
	}
	if m["out"] != nil {
		var out []string
		for _, d := range m["out"].([]interface{}) {
			out = append(out, d.(string))
		}
		b.Out = out
	}
	return &b
}

func (b Block) Run(state stateStorage.State, p *Pipeline) {
	if !b.Changed(state, p) {
		return
	}

	initState := make(map[string]string)
	for _, d := range b.DepsFull {
		var st string
		var err error
		if v, ok := p.Map[d]; ok {
			st, err = v.GetState()
		} else if v, ok := p.MapFiles[d]; ok {
			st, err = v.GetState()
		} else {
			log.Println(d, "NOT FOUND")
			continue
		}
		if err != nil {
			log.Println(err)
			st = ""
		}
		initState[d] = st
	}

	utils.GenericRun(b.Exec)
	s, err := b.GetState()
	if err != nil {
		log.Println(err)
	} else {
		state.Set(b.FullName, s)
	}

	for _, f := range b.Out {
		c, err := p.MapFiles[f].GetState()
		if err != nil {
			log.Println(err)
			continue
		}
		state.Set(f, c)
	}

	js, err := json.Marshal(initState)
	if err != nil {
		log.Println("json marshall error:", err)
		return
	}

	state.Set(utils.DepsPrefix + b.FullName, string(js))
}

func (b Block) GetState() (string, error) {
	// todo this should be used for checking external deps (i.e. download over internet)
	if len(b.Check) == 0 {
		log.Printf("phase %s does not support state checking", b.Name)
		// todo custom error
		return "", errors.New("phase does not support state check")
	}
	out := utils.GenericRun(b.Check)
	cs, err := utils.Md5Checksum(strings.NewReader(out))
	if err != nil {
		return "EMPTY", err
	}
	return cs, nil
}

func (b Block) Changed(state stateStorage.State, p *Pipeline) bool {
	// todo review
	newState, err := b.GetState()
	if err != nil {
		return true
	}
	for _, d := range b.DepsFull {
		if v, ok := p.Map[d]; ok {
			v.Changed(state, p)
		} else if v, ok := p.MapFiles[d]; ok {
			v.Changed(state, p)
		}
	}
	currentState := state.Get(b.Name)
	if currentState != "" && newState == currentState {
		log.Printf("phase %s will not run - state did not change", b.Name)
		return false
	}
	return true
}
