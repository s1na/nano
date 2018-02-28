package types

import (
	"encoding/hex"
	"encoding/json"
	"strings"
)

type Signature [64]byte

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
