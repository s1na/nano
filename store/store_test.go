package store

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	s := NewStore("unittestdir")

	s.Start()
	s.Stop()

	os.RemoveAll("unittestdir")
}

func TestSetGet(t *testing.T) {
	key := []byte("key")
	val := []byte("value")

	s := NewStore("unittestdir")

	s.Start()

	err := s.Set(key, val)
	require.Nil(t, err)

	v, err := s.Get(key)
	require.Nil(t, err)

	assert.Equal(t, val, v)

	s.Stop()

	os.RemoveAll("unittestdir")
}
