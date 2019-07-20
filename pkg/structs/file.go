package structs

import (
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"io"
	"os"
)

type File struct {
	Name  string
	Deps  []*Block
}

func (f File) openFile() (io.Reader, error) {
	file, err := os.Open(f.Name)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	return file, nil
}

func (f File) GetState() (string, error) {
	r, err := f.openFile()
	if err != nil {
		return "", err
	}
	return utils.Md5Checksum(r)
}

func (f File) Changed(state stateStorage.State, p *Pipeline) bool {
	md5Sum, err := f.GetState()
	if err != nil {
		return true
	}
	s := state.Get(f.Name)
	if s != md5Sum {
		return true
	}
	return false
}
