package blocks

import (
	"github.com/s1na/nano/types"
)

type ChangeBlock struct {
	Previous       types.BlockHash
	Representative types.PubKey
	CommonBlock
}

func (b *ChangeBlock) Hash() types.BlockHash {
	return HashChange(b.Previous, b.Representative)
}

func (b *ChangeBlock) GetPrevious() types.BlockHash {
	return b.Previous
}

func (b *ChangeBlock) GetRoot() types.BlockHash {
	return b.Previous
}

func (*ChangeBlock) Type() BlockType {
	return Change
}

func HashChange(previous types.BlockHash, representative types.PubKey) types.BlockHash {
	return HashBytes(previous[:], representative)
}
