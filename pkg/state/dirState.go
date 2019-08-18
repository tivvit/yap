package state

import (
	"encoding/json"
	"time"
)

type DirState struct {
	Exists   bool                 `json:"exists,omitempty"`
	Files    []string             `json:"files,omitempty"`
	ModTimes map[string]time.Time `json:"mod-times,omitempty"`
	Md5s     map[string]string    `json:"md5s,omitempty"`
}

func (ds DirState) Serialize() (string, error) {
	b, err := json.Marshal(ds)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (ds *DirState) Deserialize(s string) error {
	err := json.Unmarshal([]byte(s), &ds)
	return err
}
