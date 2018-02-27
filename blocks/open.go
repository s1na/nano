package blocks

import (
	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/address"
	"github.com/frankh/nano/types"
)

type OpenBlock struct {
	SourceHash     types.BlockHash
	Representative types.Account
	Account        types.Account
	CommonBlock
}

func (b *OpenBlock) Hash() types.BlockHash {
	return HashOpen(b.SourceHash, b.Representative, b.Account)
}

func (b *OpenBlock) PreviousBlockHash() types.BlockHash {
	return b.SourceHash
}

func (b *OpenBlock) RootHash() types.BlockHash {
	pub, _ := address.AddressToPub(b.Account)
	return types.BlockHashFromSlice(pub)
}

func (*OpenBlock) Type() BlockType {
	return Open
}

func (b *OpenBlock) VerifySignature() (bool, error) {
	pub, _ := address.AddressToPub(b.Account)
	return ed25519.Verify(pub, b.Hash().Slice(), b.Signature[:]), nil
}

func HashOpen(source types.BlockHash, representative types.Account, account types.Account) types.BlockHash {
	reprBytes, _ := address.AddressToPub(representative)
	accountBytes, _ := address.AddressToPub(account)

	return HashBytes(source[:], reprBytes, accountBytes)
}
