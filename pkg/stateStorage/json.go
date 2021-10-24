package stateStorage

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/storage"
)

type jsonStorage struct {
	data    map[string]string
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
	b, err := json.MarshalIndent(js.data, "", "\t")
	if err != nil {
		log.Fatalln(err)
	}
	js.storage.Write(b)
}

func NewJsonStorage(settings Settings) *jsonStorage {
	f := settings.File
	js := jsonStorage{
		storage: storage.NewFileStorage(f),
		data:    map[string]string{},
	}
	b := js.storage.Read()
	if len(b) == 0 {
		return &js
	}
	err := json.Unmarshal(b, &js.data)
	if err != nil {
		log.Fatalln(err)
	}
	return &js
}
