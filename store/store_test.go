package store

import (
	"os"
	"testing"
)

func TestStart(t *testing.T) {
	s := NewStore("unittestdir")

	s.Start()
	s.Stop()

	os.RemoveAll("unittestdir")
}
