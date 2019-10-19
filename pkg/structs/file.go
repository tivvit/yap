package structs

import (
	"github.com/emicklei/dot"
	"github.com/tivvit/yap/pkg/state"
	"github.com/tivvit/yap/pkg/stateStorage"
	"github.com/tivvit/yap/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"reflect"
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
	f.Analyze()
	if !f.Exists {
		return state.FileState{
			Exists: f.Exists,
		}.Serialize()
	}
	if f.IsDir {
		files, err := f.getDirList()
		if err != nil {
			return "", err
		}
		return state.DirState{
			Exists:   f.Exists,
			Files:    files,
			ModTimes: modTimes(files),
			Md5s:     md5s(files),
		}.Serialize()
	} else {
		md5, err := f.getFileMd5()
		if err != nil {
			return "", err
		}
		return state.FileState{
			Exists:  f.Exists,
			ModTime: f.ModTime,
			Md5:     md5,
		}.Serialize()
	}
}

func (f File) getFileMd5() (string, error) {
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

func (f File) getDirList() ([]string, error) {
	var files []string
	err := filepath.Walk(f.Name, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func modTimes(files []string) []time.Time {
	var modTimes []time.Time
	for _, f := range files {
		fo := File{
			Name: f,
		}
		fo.Analyze()
		modTimes = append(modTimes, fo.ModTime)
	}
	return modTimes
}

func md5s(files []string) []string {
	var md5s []string
	for _, f := range files {
		fo := File{
			Name: f,
		}
		md5, err := fo.getFileMd5()
		if err != nil {
			log.Println(err)
		}
		md5s = append(md5s, md5)
	}
	return md5s
}

func (f File) loadState(s string) (interface{}, error) {
	if f.IsDir {
		ds := state.DirState{}
		err := ds.Deserialize(s)
		if err != nil {
			return nil, err
		}
		return ds, nil
	} else {
		fs := state.FileState{}
		err := fs.Deserialize(s)
		if err != nil {
			return nil, err
		} else {
			return fs, nil
		}
	}
}

func (f File) Changed(s stateStorage.State, p *Pipeline) bool {
	// todo this should compare with some exact state not with my state
	oldState, err := f.loadState(s.Get(f.GetFullName()))
	log.Println("checking ", f.Name)
	if err != nil {
		// not possible to read state = recompute
		return true
	}
	f.Analyze()
	switch oldState.(type) {
	case state.FileState:
		log.Println("checking file", f.Name, oldState)
		oldFS := oldState.(state.FileState)
		if f.Exists != oldFS.Exists {
			return true
		}
		if !f.Exists {
			return false
		}
		if f.Analyzed && f.IsDir {
			return true
		}
		if f.ModTime != oldFS.ModTime {
			return true
		}
		md5, err := f.getFileMd5()
		if err != nil {
			log.Println(err)
			return true
		}
		if md5 != oldFS.Md5 {
			return true
		}
		return false
	case state.DirState:
		log.Println("checking dir", f.Name)
		oldDS := oldState.(state.DirState)
		if f.Exists != oldDS.Exists {
			return true
		}
		if !f.Exists {
			return false
		}
		if f.Analyzed && !f.IsDir {
			return true
		}
		dirList, err := f.getDirList()
		if err != nil {
			log.Println(err)
			return true
		}
		if !reflect.DeepEqual(dirList, oldDS.Files) {
			return true
		}
		modTimes := modTimes(dirList)
		if !reflect.DeepEqual(modTimes, oldDS.ModTimes) {
			return true
		}
		md5s := md5s(dirList)
		if !reflect.DeepEqual(md5s, oldDS.Md5s) {
			return true
		}
		return false
	default:
		log.Println("checking default", f.Name)
		return true
	}
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

func (f File) Visualize(ctx *dot.Graph, p *Pipeline, fileMap *map[string]*File, m *map[string]dot.Node, conf VisualizeConf) {
	if f.Analyzed && f.IsDir {
		(*m)[f.GetFullName()] = ctx.Node(DotFilePrefix+f.Name).Attr("shape", DirShape).Label(f.Name)
	} else {
		(*m)[f.GetFullName()] = ctx.Node(DotFilePrefix+f.Name).Attr("shape", FileShape).Label(f.Name)
	}
	if conf.Check && f.Changed(p.State, p) {
		log.Println("  changed")
		(*m)[f.GetFullName()].Attr("color", utils.DotChangedColor)
	}
}
