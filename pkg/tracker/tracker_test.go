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

func TestMultiTrack(t *testing.T) {
	tr := NewTracker()
	d := 5 * time.Millisecond
	tr.Start("a")
	tr.Start("b")
	time.Sleep(d)
	mda, err := tr.Stop("a")
	if err != nil {
		t.Fail()
	}
	time.Sleep(d)
	mdb, err := tr.Stop("b")
	if err != nil {
		t.Fail()
	}
	mdat := mda.Truncate(time.Millisecond)
	mdbt := mdb.Truncate(time.Millisecond)
	if mdat != d {
		t.Fail()
	}
	if mdbt != 2 * d {
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


func TestMultiRead(t *testing.T) {
	tr := NewTracker()
	tr.Start("a")
	_, err := tr.Stop("a")
	if err != nil {
		t.Fail()
	}
	_, err = tr.Stop("a")
	if err == nil {
		t.Fail()
	}
}