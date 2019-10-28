package reporter

import (
	"github.com/tivvit/yap/pkg/structs"
	"testing"
)

func TestReport(t *testing.T) {
	r := NewReporter()
	r.Report(structs.Event{
		Message: "Hi",
	})
}
