package blocks

import (
	"github.com/frankh/nano/address"
	"github.com/frankh/nano/types"
)

type ChangeBlock struct {
	PreviousHash   types.BlockHash
	Representative types.Account
	CommonBlock
}

func (b *ChangeBlock) Hash() types.BlockHash {
	return HashChange(b.PreviousHash, b.Representative)
}

func (b *ChangeBlock) PreviousBlockHash() types.BlockHash {
	return b.PreviousHash
}

func (b *ChangeBlock) RootHash() types.BlockHash {
	return b.PreviousHash
}

func (*ChangeBlock) Type() BlockType {
	return Change
}

func HashChange(previous types.BlockHash, representative types.Account) types.BlockHash {
	reprBytes, _ := address.AddressToPub(representative)
	return HashBytes(previous[:], reprBytes)
}
