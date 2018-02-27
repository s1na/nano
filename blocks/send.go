package blocks

import (
	"encoding/hex"

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
	return types.BlockHashFromBytes(HashSend(b.PreviousHash, b.Destination, b.Balance))
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

func HashSend(previous types.BlockHash, destination types.Account, balance uint128.Uint128) (result []byte) {
	previous_bytes, _ := hex.DecodeString(string(previous))
	dest_bytes, _ := address.AddressToPub(destination)
	balance_bytes := balance.GetBytes()

	return HashBytes(previous_bytes, dest_bytes, balance_bytes)
}
