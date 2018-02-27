package network

import (
	"errors"

	"github.com/golang/crypto/blake2b"
)

type Vote struct {
	Account   [32]byte
	Signature [64]byte
	Sequence  [8]byte
	Block
}

func (m *Vote) Hash() []byte {
	hash, _ := blake2b.New(32, nil)

	hash.Write(m.Block.ToBlock().Hash().Slice())
	hash.Write(m.Sequence[:])

	return hash.Sum(nil)
}

func (m *Vote) Unmarshal(data []byte) error {
	vb, bb := data[:104], data[104:]
	if len(vb) != 104 || len(bb) == 0 {
		return errors.New("invalid vote")
	}

	copy(m.Account[:], vb[:32])
	copy(m.Signature[:], vb[32:96])
	copy(m.Sequence[:], vb[96:104])

	return m.Block.Unmarshal(bb)
}

func (m *Vote) Marshal() ([]byte, error) {
	data := make([]byte, 0, 104)

	data = append(data, m.Account[:]...)
	data = append(data, m.Signature[:]...)
	data = append(data, m.Sequence[:]...)

	block, err := m.Block.Marshal()
	if err != nil {
		return nil, err
	}
	data = append(data, block...)

	return data, nil
}
