package utils

import (
	"bytes"
	"log"
	"os"
	"os/exec"
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
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(out.String())
	return out.String()
}

