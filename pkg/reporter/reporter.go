package reporter

import (
	"github.com/tivvit/yap/pkg/reporterStorage"
	"github.com/tivvit/yap/pkg/structs"
	"log"
)

type reporter struct {
	storages []reporterStorage.ReporterStorage
}

func (r *reporter) Report(e *structs.Event) {
	for _, s := range r.storages {
		s.Add(*e)
	}
}

// todo singleton
func NewReporter(rc structs.ReporterConf) *reporter {
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
	return &reporter{
		storages: storages,
	}
}
