package reporterStorage

import (
	"encoding/json"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/storage"
	"log"
)

type events []event.Event

type jsonStorage struct {
	storage storage.Storage
	events  events
}

func (js *jsonStorage) Add(e event.Event) {
	js.events = append(js.events, e)
	defer js.write()
}

func (js jsonStorage) write() {
	b, err := json.MarshalIndent(js.events, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	js.storage.Write(b)
}

func NewJsonStorage(fileName string) *jsonStorage {
	js := jsonStorage{
		storage: storage.NewFileStorage(fileName),
	}
	b := js.storage.Read()
	if len(b) == 0 {
		return &js
	}
	return &js
}
