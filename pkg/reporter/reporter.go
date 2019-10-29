package reporter

import (
	"errors"
	"github.com/tivvit/yap/pkg/reporterStorage"
	"github.com/tivvit/yap/pkg/structs"
	"log"
)

var (
	instance *reporter
)

type reporter struct {
	storages []reporterStorage.ReporterStorage
}

func (r *reporter) Report(e *structs.Event) {
	for _, s := range r.storages {
		s.Add(*e)
	}
}

func newReporter(rc structs.ReporterConf) *reporter {
	var storages []reporterStorage.ReporterStorage
	for _, s := range rc.Storages {
		switch s.(type) {
		case structs.ReporterStorageConfJson:
			rscj := s.(structs.ReporterStorageConfJson)
			// todo going to override (report?)
			storages = append(storages, reporterStorage.NewJsonStorage(rscj.FileName))
		case structs.ReporterStorageConfStdout:
			storages = append(storages, reporterStorage.NewStdoutStorage())
		default:
			log.Printf("Unknown reporter storage %T\n", s)
		}
	}
	instance =  &reporter{
		storages: storages,
	}
	return instance
}

func NewReporter(rc structs.ReporterConf) *reporter {
	if instance != nil {
		log.Println("Reporter instance already exists")
		return instance
	}
	return newReporter(rc)
}

func GetInstance() (*reporter, error) {
	if instance == nil {
		return nil, errors.New("reporter was not instantiated")
	}
	return instance, nil
}
