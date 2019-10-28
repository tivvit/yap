package reporterStorage

import (
	"github.com/tivvit/yap/pkg/structs"
)

type Storage interface {
	Add(e structs.Event)
}
