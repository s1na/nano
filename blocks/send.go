package blocks

import (
	"github.com/frankh/nano/address"
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/uint128"
)

type SendBlock struct {
	PreviousHash types.BlockHash
	Destination  types.Account
	Balance      uint128.Uint128
	CommonBlock
}

func (b *SendBlock) Hash() types.BlockHash {
	return HashSend(b.PreviousHash, b.Destination, b.Balance)
}

func (b *SendBlock) PreviousBlockHash() types.BlockHash {
	return b.PreviousHash
}

func (b *SendBlock) RootHash() types.BlockHash {
	return b.PreviousHash
}

func (*SendBlock) Type() BlockType {
	return Send
}

func HashSend(previous types.BlockHash, destination types.Account, balance uint128.Uint128) types.BlockHash {
	destBytes, _ := address.AddressToPub(destination)
	balanceBytes := balance.GetBytes()

	return HashBytes(previous[:], destBytes, balanceBytes)
}
