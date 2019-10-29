package reporter

import (
	"encoding/json"
	"github.com/tivvit/yap/pkg/structs"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fn := []string{"report.json", "report_2.json"}
	for _, f := range fn {
		os.Remove(f)
	}
	r := m.Run()
	for _, f := range fn {
		os.Remove(f)
	}
	os.Exit(r)
}


func TestStdoutReport(t *testing.T) {
	r := NewReporter(structs.ReporterConf{
		Storages: []structs.ReporterStorageConf{
			structs.ReporterStorageConfStdout{},
		},
	})
	read, w, _ := os.Pipe()
	log.SetOutput(w)
	e := structs.NewEvent()
	e.Message = "Hi"
	r.Report(e)
	w.Close()
	out, _ := ioutil.ReadAll(read)
	log.SetOutput(os.Stdout)
	if !strings.Contains(string(out), `"message": "Hi"`) {
		t.Fail()
	}
}

func TestJsonReport(t *testing.T) {
	fn := "report.json"
	r := NewReporter(structs.ReporterConf{
		Storages: []structs.ReporterStorageConf{
			structs.ReporterStorageConfJson{
				FileName: fn,
			},
		},
	})
	e := structs.NewEvent()
	e.Message = "Hi"
	r.Report(e)
	f, _ := os.Open(fn)
	b, _ := ioutil.ReadAll(f)
	if !strings.Contains(string(b), `"message": "Hi"`) {
		t.Fail()
	}

	e = structs.NewEvent()
	e.Message = "Hi 2"
	tm := time.Now()
	e.EndTime = &tm
	e.StartTime = &tm
	e.Tags = []string{"test"}
	r.Report(e)
	f, _ = os.Open(fn)
	b, _ = ioutil.ReadAll(f)
	if !strings.Contains(string(b), `"message": "Hi 2"`) {
		t.Error("message")
	}
	if !strings.Contains(string(b), `"start-time"`) {
		t.Error("start")
	}
	if !strings.Contains(string(b), `"end-time"`) {
		t.Error("end")
	}
	if !strings.Contains(string(b), `"tags": [`) {
		t.Error("tags")
	}
	var d []structs.Event
	err := json.Unmarshal(b, &d)
	if err != nil {
		t.Error("json parse")
	}
}

func TestMultiReport(t *testing.T) {
	fn := "report.json"
	fn2 := "report_2.json"
	r := NewReporter(structs.ReporterConf{
		Storages: []structs.ReporterStorageConf{
			structs.ReporterStorageConfJson{
				FileName: "report.json",
			},
			structs.ReporterStorageConfJson{
				FileName: "report_2.json",
			},
		},
	})
	e := structs.NewEvent()
	e.Message = "Hi"
	r.Report(e)
	f, _ := os.Open(fn)
	b, _ := ioutil.ReadAll(f)
	if !strings.Contains(string(b), `"message": "Hi"`) {
		t.Fail()
	}
	f2, _ := os.Open(fn2)
	b, _ = ioutil.ReadAll(f2)
	if !strings.Contains(string(b), `"message": "Hi"`) {
		t.Fail()
	}
}
