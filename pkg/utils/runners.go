package utils

import (
	"bytes"
	"github.com/tivvit/yap/pkg/reporter"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/tracker"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GenericRun(cmd []string) string {
	return run(cmd, []string{})
}

func GenericRunEnv(cmd []string, environ []string) string {
	env := os.Environ()
	for _, e := range environ {
		env = append(env, e)
	}
	return run(cmd, env)
}

func run(cmd []string, env []string) string {
	//log.Println(strings.Join(cmd, " "))
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = env
	var out bytes.Buffer
	c.Stdout = &out
	t := tracker.NewTracker()
	t.Start("run")
	err := c.Run()
	e := event.NewRunEvent("Finished")
	d, st, err := t.Stop("run")
	if err != nil {
		log.Println("Tracker for run failed")
	} else {
		e.StartTime = &st
		e.Duration = &d
	}
	e.Command = strings.Join(cmd, " ")
	e.Env = strings.Join(env, "\n")
	reporter.Report(e)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(out.String())
	return out.String()
}

