package storage

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type fileStorage struct {
	fileName string
}

func (fs fileStorage) Write(b []byte) {
	err := ioutil.WriteFile(fs.fileName, b, 0644)
	if err != nil {
		log.Warnln(err)
	}
}

func (fs fileStorage) Read() []byte {
	b, err := ioutil.ReadFile(fs.fileName)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	return b
}

func NewFileStorage(fileName string) *fileStorage {
	fs := fileStorage{
		fileName: fileName,
	}
	fs.fileName = fileName
	_, err := os.Stat(fs.fileName)
	if os.IsNotExist(err) {
		// todo report file creation
		fs.Write([]byte(""))
	}

	return &fs
}
