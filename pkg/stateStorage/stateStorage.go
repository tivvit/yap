package stateStorage

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type FileAttributes string

const (
	Mtime FileAttributes = "mtime"
	Md5                  = "md5"
)

type FileAttributeSettings struct {
	Regex       string                  `yaml:"regex"`
	Rules       map[FileAttributes]bool `yaml:"rules"`
	RegexParsed *regexp.Regexp          `yaml:"-"`
}

type Settings struct {
	Type            string `yaml:"type"`
	File            string `yaml:"file"`
	FilesAttributes []FileAttributeSettings
}

func (s *Settings) preprocessFileAttributeRegex() {
	for i, fas := range s.FilesAttributes {
		if fas.RegexParsed != nil {
			break
		}
		re, err := regexp.Compile(fas.Regex)
		if err != nil {
			log.Fatalln(err.Error())
		}
		s.FilesAttributes[i].RegexParsed = re
	}
}

func (s *Settings) FindRuleMatchingToFile(fileName string) map[FileAttributes]bool {
	s.preprocessFileAttributeRegex()
	for _, fas := range s.FilesAttributes {
		// first matching rule wins
		if fas.RegexParsed.MatchString(fileName) {
			return fas.Rules
		}
	}
	return map[FileAttributes]bool{
		Mtime: true,
		Md5:   true,
	}
}

func NewStateStorage(settings Settings) (State, error) {
	switch settings.Type {
	case "json":
		return NewJsonStorage(settings), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown storage type requested %s", settings.Type))
	}
}
