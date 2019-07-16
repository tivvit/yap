package utils

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

func GenericRun(cmd []string) string {
	log.Println(strings.Join(cmd, " "))
	c := exec.Command(cmd[0], cmd[1:]...)
	var out bytes.Buffer
	c.Stdout = &out
	err := c.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(out.String())
	return out.String()
}

