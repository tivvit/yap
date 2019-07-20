package stateStorage

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type jsonStorage map[string]string

func (js *jsonStorage) Set(key string, value string) {
	(*js)[key] = value
	js.write()
}

func (js *jsonStorage) Get(key string) string {
	return (*js)[key]
}

func (js *jsonStorage) Delete(key string) {
	delete(*js, key)
	js.write()
}

func (js jsonStorage) write() {
	b, err := json.Marshal(js)
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile("state.json", b, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewJsonStorage() *jsonStorage {
	js := jsonStorage{}
	// todo configurable
	b, err := ioutil.ReadFile("state.json")
	if err != nil {
		log.Println(err)
		return &js
	}
	err = json.Unmarshal(b, &js)
	if err != nil {
		log.Fatalln(err)
	}
	return &js
}
