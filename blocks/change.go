package blocks

import (
	"github.com/frankh/nano/types"
)

type ChangeBlock struct {
	Previous       types.BlockHash
	Representative types.AccPub
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

func HashChange(previous types.BlockHash, representative types.AccPub) types.BlockHash {
	return HashBytes(previous[:], representative)
}
