package reporter

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/reporter/event"
	"github.com/tivvit/yap/pkg/reporterStorage"
)

var (
	instance *reporter
)

type reporter struct {
	storages []reporterStorage.ReporterStorage
}

func (r *reporter) Report(e event.Event) {
	for _, s := range r.storages {
		s.Add(e)
	}
}

func newReporter(rc conf.ReporterConf) *reporter {
	var storages []reporterStorage.ReporterStorage
	for _, s := range rc.Storages {
		switch s.(type) {
		case conf.ReporterStorageConfJson:
			rscj := s.(conf.ReporterStorageConfJson)
			// todo going to override (report?)
			storages = append(storages, reporterStorage.NewJsonStorage(rscj.FileName))
		case conf.ReporterStorageConfStdout:
			storages = append(storages, reporterStorage.NewStdoutStorage())
		default:
			log.Printf("Unknown reporter storage %T\n", s)
		}
	}
	instance = &reporter{
		storages: storages,
	}
	return instance
}

func NewReporter(rc conf.ReporterConf) *reporter {
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

func Report(e event.Event) {
	r, err := GetInstance()
	if err != nil {
		log.Println("reporter instance invalid")
		return
	}
	r.Report(e)
}
