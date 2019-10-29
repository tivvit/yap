package reporterStorage

import (
	"encoding/json"
	"github.com/tivvit/yap/event"
	"log"
)

type stdoutStorage struct{}

func (js *stdoutStorage) Add(e event.Event) {
	j, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(j))
}

// todo this should take a logger instance
func NewStdoutStorage() *stdoutStorage {
	return &stdoutStorage{}
}