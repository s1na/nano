package network

import (
	"bytes"
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

	hash.Write(m.Block.ToBlock().Hash().ToBytes())
	hash.Write(m.Sequence[:])

	return hash.Sum(nil)
}

func (m *Vote) Read(messageBlockType byte, buf *bytes.Buffer) error {
	n1, err1 := buf.Read(m.Account[:])
	n2, err2 := buf.Read(m.Signature[:])
	n3, err3 := buf.Read(m.Sequence[:])

	err4 := m.Block.Read(messageBlockType, buf)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Failed to read message vote")
	}

	if n1 != 32 || n2 != 64 || n3 != 8 {
		return errors.New("Failed to read message vote")
	}

	return nil
}

func (m *Vote) Write(buf *bytes.Buffer) error {
	n1, err1 := buf.Write(m.Account[:])
	n2, err2 := buf.Write(m.Signature[:])
	n3, err3 := buf.Write(m.Sequence[:])

	err4 := m.Block.Write(buf)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return errors.New("Failed to read message vote")
	}

	if n1 != 32 || n2 != 64 || n3 != 8 {
		return errors.New("Failed to read message vote")
	}

	return nil
}
