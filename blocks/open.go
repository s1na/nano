package blocks

import (
	"github.com/frankh/crypto/ed25519"
	"github.com/s1na/nano/types"
)

type OpenBlock struct {
	Source         types.BlockHash
	Representative types.PubKey
	Account        types.PubKey
	CommonBlock
}

func (b *OpenBlock) Hash() types.BlockHash {
	return HashOpen(b.Source, b.Representative, b.Account)
}

func (b *OpenBlock) GetPrevious() types.BlockHash {
	return b.Source
}

func (b *OpenBlock) GetRoot() types.BlockHash {
	return types.BlockHashFromSlice(b.Account)
}

func (*OpenBlock) Type() BlockType {
	return Open
}

func (b *OpenBlock) VerifySignature() (bool, error) {
	return ed25519.Verify(ed25519.PublicKey(b.Account), b.Hash().Slice(), b.Signature[:]), nil
}

func HashOpen(source types.BlockHash, representative types.PubKey, account types.PubKey) types.BlockHash {
	return HashBytes(source[:], representative, account)
}
