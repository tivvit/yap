package structs

import (
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/state"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	DirShape      = "septagon"
	FileShape     = "oval"
	DotFilePrefix = "file:"
)

type File struct {
	Name     string
	Deps     []*Block
	Analyzed bool
	Size     int64
	Mode     os.FileMode
	ModTime  time.Time
	IsDir    bool
	Exists   bool
}

func (f File) GetState() (string, error) {
	if f.IsDir {
		return f.getDirState()
	} else {
		return f.getFileState()
	}
}

func (f File) getFileState() (string, error) {
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

func (f File) getDirState() (string, error) {
	var files []string
	err := filepath.Walk(f.Name, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	log.Println(files)
	//return files, err
	return "", nil
}

func (f File) changedModTime() bool {
	oldModTime := f.ModTime
	f.Analyze()
	if oldModTime != f.ModTime {
		return true
	}
	return false
}

func loadState(s string) (interface{}, error) {
	// we do not know if stored type is dir or file
	fs := state.FileState{}
	err := fs.Deserialize(s)
	if err != nil {
		// it is not fileState = trying DirState
	} else {
		return fs, nil
	}
	ds := state.DirState{}
	err = ds.Deserialize(s)
	if err != nil {
		return nil, err
	}
	return ds, nil
}

func (f File) Changed(s stateStorage.State, p *Pipeline) bool {
	oldState, err := loadState(s.Get(f.GetFullName()))
	if err != nil {
		// not possible to read state = recompute
		return true
	}
	// todo this supposes that the file has been analyzed (right before the call)
	switch oldState.(type) {
	case state.FileState:
		oldFS := oldState.(state.FileState)
		if f.Exists != oldFS.Exists {
			return true
		}
		if f.ModTime != oldFS.ModTime {
			return true
		}
		md5, err := f.getFileState()
		if err != nil {
			log.Println(err)
			return true
		}
		if md5 != oldFS.Md5 {
			return true
		}
		return false
	case state.DirState:
		oldDS := oldState.(state.DirState)
		if f.Exists != oldDS.Exists {
			return true
		}
		// todo check file list
		// todo check mod times
		// todo check md5s
	default:
		return true
	}
	return true
	//if !f.Analyzed {
	//	return true
	//}
	//if f.changedModTime() {
	//	return true
	//}
	//md5Sum, err := f.GetState()
	//if err != nil {
	//	return true
	//}
	//if s != md5Sum {
	//	return true
	//}
	//return false
}

func (f *File) Analyze() {
	fileInfo, err := os.Stat(f.Name)
	if os.IsNotExist(err) {
		f.Exists = false
	} else {
		f.Exists = true
	}
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
		r = append(r, d.GetFullName())
	}
	return r
}

func (f File) GetFullName() string {
	return DotFilePrefix + f.Name
}

func (f File) Visualize(ctx *dot.Graph, fileMap *map[string]*File, m *map[string]dot.Node, conf VisualizeConf) {
	if f.Analyzed && f.IsDir {
		(*m)[f.GetFullName()] = ctx.Node(DotFilePrefix+f.Name).Attr("shape", DirShape).Label(f.Name)
	} else {
		(*m)[f.GetFullName()] = ctx.Node(DotFilePrefix+f.Name).Attr("shape", FileShape).Label(f.Name)
	}
}
