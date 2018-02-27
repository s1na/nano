package blocks

import (
	"encoding/hex"

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
	return types.BlockHashFromBytes(HashOpen(b.SourceHash, b.Representative, b.Account))
}

func (b *OpenBlock) PreviousBlockHash() types.BlockHash {
	return b.SourceHash
}

func (b *OpenBlock) RootHash() types.BlockHash {
	pub, _ := address.AddressToPub(b.Account)
	return types.BlockHash(hex.EncodeToString(pub))
}

func (*OpenBlock) Type() BlockType {
	return Open
}

func (b *OpenBlock) VerifySignature() (bool, error) {
	pub, _ := address.AddressToPub(b.Account)
	res := ed25519.Verify(pub, b.Hash().ToBytes(), b.Signature.ToBytes())
	return res, nil
}

func HashOpen(source types.BlockHash, representative types.Account, account types.Account) (result []byte) {
	source_bytes, _ := hex.DecodeString(string(source))
	repr_bytes, _ := address.AddressToPub(representative)
	account_bytes, _ := address.AddressToPub(account)
	return HashBytes(source_bytes, repr_bytes, account_bytes)
}
