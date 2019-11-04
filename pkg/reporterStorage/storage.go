package reporterStorage

import (
	"github.com/tivvit/yap/pkg/reporter/event"
)

type ReporterStorage interface {
	Add(e event.Event)
}
