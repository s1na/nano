package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/frankh/crypto/ed25519"
)

type BlockHash [32]byte

func BlockHashFromString(s string) (BlockHash, error) {
	var h BlockHash
	r, err := hex.DecodeString(s)
	if err != nil {
		return h, err
	}

	copy(h[:], r)

	return h, nil
}

func BlockHashFromSlice(data []byte) BlockHash {
	var h BlockHash
	copy(h[:], data)

	return h
}

func (h BlockHash) IsZero() bool {
	for _, v := range h {
		if v != 0 {
			return false
		}
	}

	return true
}

func (h BlockHash) Sign(prv ed25519.PrivateKey) Signature {
	return SignatureFromSlice(ed25519.Sign(prv, h[:]))
}

func (h BlockHash) Slice() []byte {
	return h[:]
}

func (h BlockHash) String() string {
	return strings.ToUpper(hex.EncodeToString(h[:]))
}

func (h *BlockHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	r, err := BlockHashFromString(s)
	if err != nil {
		return err
	}

	*h = r

	return nil
}

func (h BlockHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}
