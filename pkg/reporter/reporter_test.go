package reporter

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tivvit/yap/pkg/conf"
	"github.com/tivvit/yap/pkg/reporter/event"
	"io/ioutil"
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
	r := newReporter(conf.ReporterConf{
		Storages: []conf.ReporterStorageConf{
			conf.ReporterStorageConfStdout{},
		},
	})
	read, w, _ := os.Pipe()
	log.SetOutput(w)
	e := event.NewEvent("Hi")
	r.Report(e)
	w.Close()
	out, _ := ioutil.ReadAll(read)
	log.SetOutput(os.Stdout)
	//log.Println(string(out))
	// todo this output is not nice - when using logrus use structured logs
	if !strings.Contains(string(out), `\"message\": \"Hi\"`) {
		t.Fail()
	}
}

func TestJsonReport(t *testing.T) {
	fn := "report.json"
	r := newReporter(conf.ReporterConf{
		Storages: []conf.ReporterStorageConf{
			conf.ReporterStorageConfJson{
				FileName: fn,
			},
		},
	})
	e := event.NewEvent("Hi")
	r.Report(e)
	f, _ := os.Open(fn)
	b, _ := ioutil.ReadAll(f)
	if !strings.Contains(string(b), `"message": "Hi"`) {
		t.Fail()
	}

	bre := event.NewBlockRunEvent("Hi 2", "block")
	tm := time.Now()
	bre.StartTime = &tm
	bre.Tags = []string{"test"}
	r.Report(bre)
	f, _ = os.Open(fn)
	b, _ = ioutil.ReadAll(f)
	if !strings.Contains(string(b), `"message": "Hi 2"`) {
		t.Error("message")
	}
	if !strings.Contains(string(b), `"start-time"`) {
		t.Error("start")
	}
	if !strings.Contains(string(b), `"tags": [`) {
		t.Error("tags")
	}
	if !strings.Contains(string(b), `"block": "block"`) {
		t.Error("block")
	}
	var d []event.BaseEvent
	err := json.Unmarshal(b, &d)
	if err != nil {
		t.Error("json parse")
	}
}

func TestMultiReport(t *testing.T) {
	fn := "report.json"
	fn2 := "report_2.json"
	r := newReporter(conf.ReporterConf{
		Storages: []conf.ReporterStorageConf{
			conf.ReporterStorageConfJson{
				FileName: "report.json",
			},
			conf.ReporterStorageConfJson{
				FileName: "report_2.json",
			},
		},
	})
	e := event.NewEvent("Hi")
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

func TestInstance(t *testing.T) {
	instance = nil
	_, err := GetInstance()
	if err == nil {
		t.Error("got uninitialized instance")
	}
	ri := newReporter(conf.ReporterConf{
		Storages: []conf.ReporterStorageConf{
			conf.ReporterStorageConfStdout{},
		},
	})
	r, err := GetInstance()
	if err != nil {
		t.Error("error getting instance")
	}
	if r == nil {
		t.Error("got nil instance")
	}
	if ri != r {
		t.Error("instances differ")
	}
	ri2 := NewReporter(conf.ReporterConf{
		Storages: []conf.ReporterStorageConf{
			conf.ReporterStorageConfStdout{},
		},
	})
	if ri2 != r {
		t.Error("instances differ after second init")
	}
}
