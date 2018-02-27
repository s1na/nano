package blocks

import (
	"github.com/frankh/nano/types"
)

type ReceiveBlock struct {
	PreviousHash types.BlockHash
	SourceHash   types.BlockHash
	CommonBlock
}

func (b *ReceiveBlock) Hash() types.BlockHash {
	return HashReceive(b.PreviousHash, b.SourceHash)
}

func (b *ReceiveBlock) PreviousBlockHash() types.BlockHash {
	return b.PreviousHash
}

func (b *ReceiveBlock) RootHash() types.BlockHash {
	return b.PreviousHash
}

func (*ReceiveBlock) Type() BlockType {
	return Receive
}

func HashReceive(previous types.BlockHash, source types.BlockHash) types.BlockHash {
	return HashBytes(previous[:], source[:])
}
