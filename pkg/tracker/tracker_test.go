package tracker

import (
	"testing"
	"time"
)

func TestTrack(t *testing.T) {
	tr := NewTracker()
	tr.Start("a")
	d := 5 * time.Millisecond
	time.Sleep(d)
	md, err := tr.Stop("a")
	if err != nil {
		t.Fail()
	}
	mdt := md.Truncate(time.Millisecond)
	if mdt != d {
		t.Fail()
	}
}


func TestMissing(t *testing.T) {
	tr := NewTracker()
	_, err := tr.Stop("a")
	if err == nil {
		t.Fail()
	}
}