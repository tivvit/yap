package structs

import (
	"encoding/json"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/emicklei/dot"
	"github.com/mattn/go-shellwords"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/conf"
	yapDot "github.com/tivvit/yap/pkg/dot"
	"github.com/tivvit/yap/pkg/reporter"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/tracker"
	"github.com/tivvit/yap/pkg/utils"
	"strconv"
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
	Stdout            bool     `yaml:"stdout,omitempty"`
	Stderr            bool     `yaml:"stderr,omitempty"`
	MayFail           bool     `yaml:"may-fail,omitempty"`
	Idempotent        bool     `yaml:"idempotent,omitempty"`
}

type IncludeBlock struct {
	Include string
	Type    string
}

const (
	dotBlockPrefix        = "block:"
	StateNameExec         = "exec"
	StateNameEnv          = "env"
	StateNameCheckCmd     = "checkCmd"
	StateNameCheck        = "checkOutHash"
	StateNameCheckSuccess = "checkOutSuccess"
)

func NewBlockFromMap(name string, m map[string]interface{}) *Block {
	b := Block{
		PipelineBlockBase: PipelineBlockBase{},
	}
	b.Name = name
	b.Stdout = true
	b.Stderr = true
	b.Idempotent = true
	if m["stdout"] != nil {
		b.Stdout = m["stdout"].(bool)
	}
	if m["stderr"] != nil {
		b.Stderr = m["stderr"].(bool)
	}
	if m["may-fail"] != nil {
		b.MayFail = m["may-fail"].(bool)
	}
	if m["idempotent"] != nil {
		b.Idempotent = m["idempotent"].(bool)
	}
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
		for _, i := range m["in"].([]interface{}) {
			b.In = append(b.In, filesFoundOrListed(i.(string))...)
		}
	}
	if m["out"] != nil {
		for _, d := range m["out"].([]interface{}) {
			b.Out = append(b.Out, filesFoundOrListed(d.(string))...)
		}
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

func filesFoundOrListed(filePattern string) []string {
	var files []string
	g, err := doublestar.Glob(filePattern)
	if err != nil {
		log.Warnf("Glob %s failed: %s", filePattern, g)
	} else {
		if len(g) > 0 {
			for _, f := range g {
				files = append(files, f)
			}
		} else {
			// files do not exist yet - we have to use what user provided
			files = append(files, filePattern)
		}
	}
	return files
}

func (b *Block) Run(state stateStorage.State, p *Pipeline, dry bool) {
	t := tracker.NewTracker()
	dryLog := ""
	if dry {
		dryLog = "dry "
	}
	log.Infof("%s running %s `%s`", dryLog, b.FullName, strings.Join(b.Exec, " "))

	if !b.Changed(state, p) {
		e := event.NewBlockRunEvent("Not changed", b.FullName)
		if !dry {
			reporter.Report(e)
		}
		return
	}

	if dry {
		return
	}

	initState := b.GetDepsState(p)
	t.Start(b.FullName)
	_, ok := utils.GenericRunEnv(b.Exec, b.Env, b.Stdout, b.Stderr)
	fatalFail := false
	if !ok {
		log.Warnf("%s failed", b.FullName)
		if !b.MayFail {
			fatalFail = true
		}
	}
	d, st, err := t.Stop(b.FullName)
	e := event.NewBlockRunEvent("Finished", b.FullName)
	if err != nil {
		log.Printf("invalid tracker for %s", b.FullName)
	} else {
		e.StartTime = &st
		e.Duration = &d
	}
	reporter.Report(e)

	// update state after run
	s, err := b.GetState()
	if err != nil {
		log.Println(err)
	} else {
		state.Set(b.FullName, s)
	}

	// state of input files has to be stored somewhere - this is probably the correct place because the block is using the files
	for _, f := range append(b.Out, b.In...) {
		file := p.FilesMap[f]
		c, err := file.GetState()
		if err != nil {
			log.Println(err)
			continue
		}
		state.Set(file.GetFullName(), c)
		e := event.NewBlockRunEvent(fmt.Sprintf("Setting state for %s", f), b.FullName)
		e.AddTag("state")
		reporter.Report(e)
	}

	log.Infof("Finished %s in %s", b.FullName, d)
	if fatalFail {
		log.Fatalf("%s failed - stopping pipeline", b.FullName)
		// todo this should be passed higher and pipeline should end correctly
	}
	state.Set(utils.DepsPrefix+b.FullName, initState)
}

func (b Block) GetDepsState(p *Pipeline) string {
	depsState := make(map[string]string)
	for _, fp := range b.In {
		g, err := doublestar.Glob(fp)
		if err != nil {
			log.Warnf("Glob %s failed: %s", fp, g)
		} else {
			for _, f := range g {
				var st string
				var err error
				if v, ok := p.FilesMap[f]; ok {
					st, err = v.GetState()
				} else {
					log.Println(f, "NOT FOUND")
					continue
				}
				if err != nil {
					log.Println(err)
					st = ""
				}
				depsState[f] = st
			}
		}
	}
	for _, d := range b.DepsFull {
		var st string
		var err error
		if v, ok := p.BlocksMap[d]; ok {
			st, err = v.GetState()
		} else {
			log.Println(d, "NOT FOUND")
			continue
		}
		if err != nil {
			log.Println(err)
			st = ""
		}
		depsState[d] = st
	}
	js, err := json.Marshal(depsState)
	if err != nil {
		log.Println("json marshall error:", err)
		return "{}"
	}
	return string(js)
}

func (b Block) GetState() (string, error) {
	state := map[string]string{}
	// add command to state
	state[StateNameExec] = strings.Join(b.Exec, " ")
	// add env to state
	state[StateNameEnv] = strings.Join(b.Env, ",")
	// add check command to state
	state[StateNameCheckCmd] = strings.Join(b.Check, " ")

	if len(b.Check) > 0 {
		// this should be used for checking external deps (i.e. download over internet)
		log.Printf("Explicit state check `%s`\n", strings.Join(b.Check, " "))
		out, ok := utils.SilentRun(b.Check)
		cs, err := utils.Md5Checksum(strings.NewReader(out))
		if err != nil {
			log.Println("Check out md5 err:", err)
			state[StateNameCheck] = "ERROR"
		}
		state[StateNameCheck] = cs
		state[StateNameCheckSuccess] = strconv.FormatBool(ok)
	}

	js, err := json.Marshal(state)
	if err != nil {
		log.Println("json marshall error:", err)
		return "{}", err
	}
	return string(js), nil
}

func (b Block) Changed(state stateStorage.State, p *Pipeline) bool {
	if !b.Idempotent {
		return true
	}
	currentState, err := b.GetState()
	if err != nil {
		log.Printf("Error while checking state for %s\n", b.FullName)
		return true
	}
	// json should have persistent ordering of keys in json
	storedState := state.Get(b.FullName)
	if storedState != currentState {
		log.Printf("%s state changed", b.FullName)
		return true
	}

	// check deps
	depsState := b.GetDepsState(p)
	storedDepsState := state.Get(utils.DepsPrefix + b.FullName)
	if depsState != storedDepsState {
		log.Printf("%s deps state changed", b.FullName)
		return true
	}
	log.Printf("not running %s - state did not change", b.Name)
	return false
}

func (b Block) GetDepsFull() []string {
	ret := make([]string, len(b.DepsFull))
	copy(ret, b.DepsFull)
	for _, f := range b.In {
		ret = append(ret, DotFilePrefix+f)
	}
	return ret
}

func (b Block) Visualize(ctx *dot.Graph, p *Pipeline, fileMap *map[string]*File, m *map[string]dot.Node, conf conf.VisualizeConf) {
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
	n := ctx.Node(dotBlockPrefix+b.FullName).Attr("shape", "plain")
	n.Attr("label", dot.HTML(label))
	if conf.Check && b.Changed(p.State, p) {
		n.Attr("color", utils.DotChangedColor)
	}
	(*m)[b.FullName] = n

	for _, f := range b.In {
		(*fileMap)[f] = nil
	}
	for _, f := range b.Out {
		(*fileMap)[f] = nil
	}
}
