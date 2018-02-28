package blocks

import (
	"github.com/frankh/nano/types"
)

type ReceiveBlock struct {
	Previous types.BlockHash
	Source   types.BlockHash
	CommonBlock
}

func (b *ReceiveBlock) Hash() types.BlockHash {
	return HashReceive(b.Previous, b.Source)
}

func (b *ReceiveBlock) GetPrevious() types.BlockHash {
	return b.Previous
}

func (b *ReceiveBlock) GetRoot() types.BlockHash {
	return b.Previous
}

func (*ReceiveBlock) Type() BlockType {
	return Receive
}

func HashReceive(previous types.BlockHash, source types.BlockHash) types.BlockHash {
	return HashBytes(previous[:], source[:])
}
