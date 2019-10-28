package reporterStorage

import (
	"encoding/json"
	"github.com/tivvit/yap/pkg/structs"
	"io/ioutil"
	"log"
	"os"
)

type events []structs.Event

type jsonStorage struct {
	fileName string
	events events
}

func (js *jsonStorage) Add(e structs.Event) {
	js.events = append(js.events, e)
	defer js.write()
}

func (js jsonStorage) write() {
	b, err := json.MarshalIndent(js.events, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile(js.fileName, b, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewJsonStorage(fileName string) *jsonStorage {
	js := jsonStorage{}
	js.fileName = fileName
	_, err := os.Stat(js.fileName)
	if os.IsNotExist(err) {
		js.write()
	}
	b, err := ioutil.ReadFile(js.fileName)
	if err != nil {
		log.Println(err)
		return &js
	}
	err = json.Unmarshal(b, &js.events)
	if err != nil {
		log.Fatalln(err)
	}
	return &js
}