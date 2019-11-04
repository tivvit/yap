package storage

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fn := "test.txt"
	os.Remove(fn)
	r := m.Run()
	os.Remove(fn)
	os.Exit(r)
}

func TestCreateReadWrite(t *testing.T) {
	fn := "test.txt"
	fs := NewFileStorage(fn)
	b := fs.Read()
	if string(b) != "" {
		t.Fail()
	}
	fs.Write([]byte("abc"))
	b = fs.Read()
	if string(b) != "abc" {
		t.Fail()
	}
	fs.Write([]byte("def"))
	b = fs.Read()
	if string(b) != "def" {
		t.Fail()
	}
}
