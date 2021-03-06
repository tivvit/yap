package reporterStorage

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/reporter/event"
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
