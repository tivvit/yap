package utils

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/reporter"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/tracker"
	"io"

	//"io"
	"os"
	"os/exec"
	"strings"
)

func GenericRun(cmd []string) (string, bool) {
	return run(cmd, []string{}, true, true)
}

func SilentRun(cmd []string) (string, bool) {
	return run(cmd, []string{}, false, false)
}

func GenericRunEnv(cmd []string, environ []string, stdout bool, stderr bool) (string, bool) {
	env := os.Environ()
	for _, e := range environ {
		env = append(env, e)
	}
	return run(cmd, env, stdout, stderr)
}

func run(cmd []string, env []string, stdout bool, stderr bool) (string, bool) {
	//log.Println(strings.Join(cmd, " "))
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = env
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	t := tracker.NewTracker()
	t.Start("run")
	if stdout {
		c.Stdout = io.MultiWriter(&out, os.Stdout)
	}
	if stderr {
		c.Stderr = io.MultiWriter(&out, os.Stderr)
	}
	err := c.Run()
	if err != nil {
		log.Warnln(err)
	}
	// todo read info about process (cpu time ...)
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
	e.Failed  = !c.ProcessState.Success()
	reporter.Report(e)
	log.Info(out.String())
	return out.String(), c.ProcessState.Success()
}
