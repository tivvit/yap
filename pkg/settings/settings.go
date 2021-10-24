package settings

import "github.com/tivvit/yap/pkg/stateStorage"

type Settings struct {
	State stateStorage.Settings `yaml:"state"`
}

func DefaultSettings() *Settings {
	return &Settings{
		State: stateStorage.Settings{
			Type: "json",
			File: "state.json",
		},
	}
}