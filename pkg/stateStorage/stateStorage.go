package stateStorage

import (
	"errors"
	"fmt"
)

type Settings struct {
	Type string `yaml:"type"`
	File string `yaml:"file"`
}

func NewStateStorage(settings Settings) (State, error) {
	switch settings.Type {
	case "json":
		return NewJsonStorage(settings), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown storage type requested %s", settings.Type))
	}
}
