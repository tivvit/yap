package reporter

import (
	"github.com/tivvit/yap/pkg/reporterStorage"
	"github.com/tivvit/yap/pkg/structs"
)

type reporter struct {
	storage reporterStorage.Storage
}

func (r *reporter) Report(e structs.Event) {
	r.storage.Add(e)
}

// todo singleton
func NewReporter() *reporter {
	// todo determine storage based on conf
	s := reporterStorage.NewStdoutStorage()
	return &reporter{
		storage: s,
	}
}
