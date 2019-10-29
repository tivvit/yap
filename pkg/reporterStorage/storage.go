package reporterStorage

import (
	"github.com/tivvit/yap/event"
)

type ReporterStorage interface {
	Add(e event.Event)
}
