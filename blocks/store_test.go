package blocks

import (
	"os"
	"testing"

	"github.com/s1na/nano/store"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetGetGenesis(t *testing.T) {
	GenesisBlock = TestGenesisBlock

	s := store.NewStore("testdata")
	s.Start()
	bs := NewBlockStore(s)

	err := bs.SetBlock(GenesisBlock)
	require.Nil(t, err)

	b, err := bs.GetBlock(GenesisBlock.Hash())
	require.Nil(t, err)

	ob := b.(*OpenBlock)
	assert.Equal(t, GenesisBlock, ob)

	s.Stop()
	os.RemoveAll("testdata")
}
