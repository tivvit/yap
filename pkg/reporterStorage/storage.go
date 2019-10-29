package reporterStorage

import (
	"github.com/tivvit/yap/pkg/structs"
)

type ReporterStorage interface {
	Add(e structs.Event)
}
