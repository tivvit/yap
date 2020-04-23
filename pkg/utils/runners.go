package utils

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/reporter"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/tracker"
	"github.com/creack/pty"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
	c := exec.Command(cmd[0], cmd[1:]...)
	c.Env = env
	out := bytes.Buffer{}
	t := tracker.NewTracker()
	t.Start("run")
	if stdout {
		stdoutPty, tty, err := pty.Open()
		if err != nil {
			log.Fatal("Not possible to open pty")
		}
		defer func() { _ = tty.Close() }()
		defer func() { _ = stdoutPty.Close() }()

		c.Stdout = tty
		if c.SysProcAttr == nil {
			c.SysProcAttr = &syscall.SysProcAttr{}
		}
		c.SysProcAttr.Setctty = true
		c.SysProcAttr.Setsid = true
		c.SysProcAttr.Ctty = int(tty.Fd())
		go func() {
			_, _ = io.Copy(io.MultiWriter(os.Stdout, &out), stdoutPty)
		}()
	} else {
		c.Stdout = io.MultiWriter(&out)
	}
	if stderr {
		stderrPty, etty, err := pty.Open()
		if err != nil {
			log.Fatal("Not possible to open pty")
		}
		defer func() { _ = etty.Close() }()
		defer func() { _ = stderrPty.Close() }()

		c.Stderr = etty
		if !stdout {
			if c.SysProcAttr == nil {
				c.SysProcAttr = &syscall.SysProcAttr{}
			}
			c.SysProcAttr.Setctty = true
			c.SysProcAttr.Setsid = true
			c.SysProcAttr.Ctty = int(etty.Fd())
		}
		go func() {
			_, _ = io.Copy(io.MultiWriter(os.Stderr, &out), stderrPty)
		}()
	} else {
		c.Stderr = io.MultiWriter(&out)
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
	e.Failed = !c.ProcessState.Success()
	reporter.Report(e)
	return out.String(), c.ProcessState.Success()
}
