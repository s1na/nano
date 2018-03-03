package blocks

import (
	"github.com/s1na/nano/types"
	"github.com/s1na/nano/uint128"

	"github.com/frankh/crypto/ed25519"
)

type UtxBlock struct {
	Account        types.PubKey
	Previous       types.BlockHash
	Representative types.PubKey
	Balance        uint128.Uint128
	Amount         uint128.Uint128
	Link           types.PubKey
	CommonBlock
}

func (b *UtxBlock) Hash() types.BlockHash {
	// TODO: Add type as uint256 preamble
	return HashUtx(b.Account, b.Previous, b.Representative, b.Balance, b.Amount, b.Link)
}

func (b *UtxBlock) GetPrevious() types.BlockHash {
	return b.Previous
}

func (b *UtxBlock) GetRoot() types.BlockHash {
	root := b.Previous
	if root.IsZero() {
		root = types.BlockHashFromSlice(b.Account[:])
	}

	return root
}

func (b *UtxBlock) Type() BlockType {
	return Utx
}

func (b *UtxBlock) VerifySignature() (bool, error) {
	return ed25519.Verify(ed25519.PublicKey(b.Account), b.Hash().Slice(), b.Signature[:]), nil
}

func (b *UtxBlock) IsSend() bool {
	// TODO: (amount.bytes [0] & 0x80) == 0x80;
	return b.Amount.Hi != 0 && b.Amount.Lo != 0
}

func HashUtx(account types.PubKey, prev types.BlockHash, repr types.PubKey, balance uint128.Uint128, amount uint128.Uint128, link types.PubKey) types.BlockHash {
	return HashBytes(account, prev[:], repr, balance.GetBytes(), amount.GetBytes(), link)
}
