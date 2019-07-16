package structs

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

type File struct {
	Name  string
	Deps  []*Block
}

func (f File) Checksum() (string, error) {
	file, err := os.Open(f.Name)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes), nil

}
func (f File) Changed(state State, p *Pipeline) bool {
	md5Sum, err := f.Checksum()
	if err != nil {
		return true
	}
	s := state.Get(f.Name)
	if s != md5Sum {
		return true
	}
	return false
}
