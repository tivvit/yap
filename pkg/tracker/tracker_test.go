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
	t.Log(md, mdt, mdt - d)
	if (mdt - d) > (2 * time.Millisecond) {
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
	t.Log(mda, mdb, mdat, mdbt, mdat -  d, mdbt - (2 * d))
	if (mdat -  d) > (2 * time.Millisecond) {
		t.Fail()
	}
	if (mdbt - (2 * d)) > (4 * time.Millisecond) {
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