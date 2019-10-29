package stateStorage

import (
	"encoding/json"
	"github.com/tivvit/yap/pkg/storage"
	"log"
)

type jsonStorage struct {
	data map[string]string
	storage storage.Storage
}

func (js *jsonStorage) Set(key string, value string) {
	(*js).data[key] = value
	js.write()
}

func (js *jsonStorage) Get(key string) string {
	return (*js).data[key]
}

func (js *jsonStorage) Delete(key string) {
	delete((*js).data, key)
	js.write()
}

func (js jsonStorage) write() {
	b, err := json.MarshalIndent(js, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	js.storage.Write(b)
}

func NewJsonStorage() *jsonStorage {
	// todo configurable
	f := "state.json"
	js := jsonStorage{
		storage: storage.NewFileStorage(f),
	}
	b := js.storage.Read()
	if len(b) == 0 {
		return &js
	}
	err := json.Unmarshal(b, &js)
	if err != nil {
		log.Fatalln(err)
	}
	return &js
}
