package blocks

import (
	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/types"
)

type UtxBlock struct {
	Account        types.AccPub
	Previous       types.BlockHash
	Representative types.AccPub
	Balance        uint128.Uint128
	Amount         uint128.Uint128
	Link           [64]byte
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
	if root.isZero() {
		root = b.Account
	}

	return types.BlockHashFromSlice(root)
}

func (b *UtxBlock) Type() BlockType {
	return Utx
}

func (b *UtxBlock) VerifySignature() (bool, error) {
	return ed25519.Verify(ed25519.PublicKey(b.Account), b.Hash().Slice(), b.Signature[:]), nil
}

func (b *UtxBlock) IsSend() bool {
	// TODO: (amount.bytes [0] & 0x80) == 0x80;
	return b.Amount != 0
}

func HashUtx(account types.AccPub, prev types.BlockHash, repr types.AccPub, balance uint128.Uint128, amount uint128.Uint128, link types.AccPub) types.BlockHash {
	return HashBytes(account, prev, repr, balance, amount, link)
}
