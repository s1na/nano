package blocks

import (
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/uint128"
)

type SendBlock struct {
	Previous    types.BlockHash
	Destination types.AccPub
	Balance     uint128.Uint128
	CommonBlock
}

func (b *SendBlock) Hash() types.BlockHash {
	return HashSend(b.Previous, b.Destination, b.Balance)
}

func (b *SendBlock) GetPrevious() types.BlockHash {
	return b.Previous
}

func (b *SendBlock) GetRoot() types.BlockHash {
	return b.Previous
}

func (*SendBlock) Type() BlockType {
	return Send
}

func HashSend(previous types.BlockHash, destination types.AccPub, balance uint128.Uint128) types.BlockHash {
	balanceBytes := balance.GetBytes()
	return HashBytes(previous[:], destination, balanceBytes)
}
