package state

import (
	"encoding/json"
	"time"
)

type FileState struct {
	Exists  bool      `json:"exists,omitempty"`
	ModTime time.Time `json:"mod-time,omitempty"`
	Md5     string    `json:"md5,omitempty"`
}

func (fs FileState) Serialize() (string, error) {
	b, err := json.Marshal(fs)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (fs *FileState) Deserialize(s string) error {
	err := json.Unmarshal([]byte(s), &fs)
	return err
}
