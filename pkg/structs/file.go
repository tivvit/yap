package structs

import (
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"log"
	"os"
	"time"
)

const (
	DirShape  = "septagon"
	FileShape = "oval"
)

type File struct {
	Name     string
	Deps     []*Block
	Analyzed bool
	Size     int64
	Mode     os.FileMode
	ModTime  time.Time
	IsDir    bool
}

func (f File) GetState() (string, error) {
	r, err := os.Open(f.Name)
	if err != nil {
		return "", err
	}

	defer func() {
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

func (f *File) Analyze() {
	fileInfo, err := os.Stat(f.Name)
	if err != nil {
		f.Analyzed = false
		log.Println(err)
		return
	}
	f.Analyzed = true
	f.Size = fileInfo.Size()
	f.Mode = fileInfo.Mode()
	f.ModTime = fileInfo.ModTime()
	f.IsDir = fileInfo.IsDir()
}

func (f File) GetDepsFull() []string {
	var r []string
	for _, d := range f.Deps {
		r = append(r, d.FullName)
	}
	return r
}

func (f File) GetFullName() string {
	return f.Name
}

func (f File) Visualize(ctx *dot.Graph, fileMap *map[string]*File, m *map[string]dot.Node, conf VisualizeConf) {
	if f.Analyzed && f.IsDir {
		(*m)[f.Name] = ctx.Node(f.Name).Attr("shape", DirShape)
	} else {
		(*m)[f.Name] = ctx.Node(f.Name).Attr("shape", FileShape)
	}
}
