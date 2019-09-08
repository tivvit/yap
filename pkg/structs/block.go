package structs

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/dot"
	"github.com/mattn/go-shellwords"
	yapDot "github.com/tivvit/yap/pkg/dot"
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
	In                []string `yaml:"in,omitempty"`
	Out               []string `yaml:"out,omitempty"`
	Env               []string `yaml:"env,omitempty"`
}

type IncludeBlock struct {
	Include string
	Type    string
}

const (
	dotBlockPrefix = "block:"
	StateNamePrefix = "__internal:"
	StateNameExec = "exec"
	StateNameEnv = "env"
)

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
	if m["in"] != nil {
		var in []string
		for _, i := range m["in"].([]interface{}) {
			in = append(in, i.(string))
		}
		b.In = in
	}
	if m["out"] != nil {
		var out []string
		for _, d := range m["out"].([]interface{}) {
			out = append(out, d.(string))
		}
		b.Out = out
	}
	if m["env"] != nil {
		var env []string
		for _, d := range m["env"].([]interface{}) {
			env = append(env, d.(string))
		}
		b.Env = env
	}
	return &b
}

func (b Block) Run(state stateStorage.State, p *Pipeline) {
	if !b.Changed(state, p) {
		return
	}

	initState := b.GetDepsState(p)

	utils.GenericRunEnv(b.Exec, b.Env)
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

	state.Set(utils.DepsPrefix+b.FullName, string(js))
}

func (b Block) GetDepsState(p *Pipeline) map[string]string {
	initState := make(map[string]string)
	for _, f := range b.In {
		var st string
		var err error
		if v, ok := p.MapFiles[f]; ok {
			st, err = v.GetState()
		} else {
			log.Println(f, "NOT FOUND")
			continue
		}
		if err != nil {
			log.Println(err)
			st = ""
		}
		initState[f] = st
	}
	for _, d := range b.DepsFull {
		var st string
		var err error
		if v, ok := p.Map[d]; ok {
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
	// add command to state
	initState[StateNamePrefix+StateNameExec] = strings.Join(b.Exec, ",")
	initState[StateNamePrefix+StateNameEnv] = strings.Join(b.Env, ",")
	return initState
}

func (b Block) GetState() (string, error) {
	// todo this should be used for checking external deps (i.e. download over internet)
	if len(b.Check) == 0 {
		// no explicit state check
		// todo get state with deps? - probably not - this may serve as interface for others to get state of this block
		return "", nil
	}
	log.Printf("Explicit state check `%s`\n", strings.Join(b.Check, " "))
	out := utils.GenericRun(b.Check)
	cs, err := utils.Md5Checksum(strings.NewReader(out))
	if err != nil {
		return "EMPTY", err
	}
	return cs, nil
}

func (b Block) Changed(state stateStorage.State, p *Pipeline) bool {
	// todo review - should be based on deps state
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

func (b Block) GetDepsFull() []string {
	ret := make([]string, len(b.DepsFull))
	copy(ret, b.DepsFull)
	for _, f := range b.In {
		ret = append(ret, DotFilePrefix + f)
	}
	return ret
}

func (b Block) Visualize(ctx *dot.Graph, fileMap *map[string]*File, m *map[string]dot.Node , conf VisualizeConf) {
	nameFmt := "<tr><td><b>%s</b></td></tr>"
	name := fmt.Sprintf(nameFmt, b.Name)
	cmdFmt := `<tr><td><font face="Courier New, Courier, monospace">%s</font></td></tr>`
	cmd := fmt.Sprintf(cmdFmt, yapDot.EscapeHtml(strings.Join(b.Exec, " ")))
	descFmt := "<tr><td>%s</td></tr>"
	desc := ""
	if b.Description != "" {
		desc = fmt.Sprintf(descFmt, b.Description)
	}
	tableFmt := `<table border="0" cellborder="1" cellspacing="0">%s%s%s</table>`
	label := fmt.Sprintf(fmt.Sprintf(tableFmt, name, desc, cmd))
	n := ctx.Node(dotBlockPrefix + b.FullName).Attr("shape", "plain")
	n.Attr("label", dot.HTML(label))
	(*m)[b.FullName] = n

	for _, f := range b.In {
		(*fileMap)[f] = nil
	}
	for _, f := range b.Out {
		(*fileMap)[f] = nil
	}
}
