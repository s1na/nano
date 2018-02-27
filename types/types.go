package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/frankh/crypto/ed25519"
)

type BlockHash [32]byte
type Account string
type Work string
type Signature [64]byte

func NewBlockHash(s string) (BlockHash, error) {
	var h BlockHash
	b, err := hex.DecodeString(s)
	if err != nil {
		return h, err
	}

	copy(h[:], b)

	return h, nil
}

func BlockHashFromSlice(b []byte) BlockHash {
	var h BlockHash
	copy(h[:], b)

	return h
}

func (h BlockHash) Slice() []byte {
	return h[:]
}

func (h BlockHash) Sign(prv ed25519.PrivateKey) Signature {
	return SignatureFromSlice(ed25519.Sign(prv, h[:]))
}

func (h BlockHash) String() string {
	return strings.ToUpper(hex.EncodeToString(h[:]))
}

func (h *BlockHash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	b, err := NewBlockHash(s)
	if err != nil {
		return err
	}

	*h = b

	return nil
}

func (h BlockHash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func NewSignature(str string) (Signature, error) {
	var s Signature
	b, err := hex.DecodeString(str)
	if err != nil {
		return s, err
	}

	copy(s[:], b)

	return s, nil
}

func SignatureFromSlice(data []byte) Signature {
	var s Signature
	copy(s[:], data)
	return s
}

func (s Signature) Slice() []byte {
	return s[:]
}

func (s Signature) String() string {
	return strings.ToUpper(hex.EncodeToString(s[:]))
}

func (s *Signature) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	b, err := NewSignature(str)
	if err != nil {
		return err
	}

	*s = b

	return nil
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
