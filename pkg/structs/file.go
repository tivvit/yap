package structs

import (
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"log"
	"os"
)

type File struct {
	Name  string
	Deps  []*Block
	// todo isDir
	// todo create time
}

func (f File) GetState() (string, error) {
	r, err := os.Open(f.Name)
	if err != nil {
		return "", err
	}

	defer func () {
		err := r.Close()
		if err != nil {
			log.Println(err)
		}
	}()

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
