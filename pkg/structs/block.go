package structs

import (
	"errors"
	"github.com/mattn/go-shellwords"
	"github.com/tivvit/yap/pkg/utils"
	"log"
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

func (b Block) Run(state State) {
	// todo resolve whole name (with path)
	if !b.Changed(state) {
		return
	}
	utils.GenericRun(b.Exec)
	s, err := b.checkState()
	if err == nil {
		state.Set(b.Name, s)
	}
	// todo change state of all files
}

func (b Block) checkState() (string, error) {
	// todo support checking based on output? what is the check for then?
	if len(b.Check) == 0 {
		log.Printf("phase %s does not support state checking", b.Name)
		// todo custom error
		return "", errors.New("phase does not support state check")
	}
	return utils.GenericRun(b.Check), nil
}

func (b Block) Changed(state State) bool {
	newState, err := b.checkState()
	if err != nil {
		return true
	}
	// todo extract file deps
	//for _, i := range b.DepsFile {}
	currentState := state.Get(b.Name)
	if currentState != "" && newState == currentState {
		log.Printf("phase %s will not run - state did not change", b.Name)
		return false
	}
	return true
}
