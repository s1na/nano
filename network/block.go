package network

import (
	"errors"

	"github.com/frankh/nano/blocks"
	"github.com/frankh/nano/types"
	"github.com/frankh/nano/uint128"
)

const (
	sendSize    = 32 + 32 + 16 + 64 + 8
	openSize    = 32 + 32 + 32 + 64 + 8
	changeSize  = 32 + 32 + 64 + 8
	receiveSize = 32 + 32 + 64 + 8
	utxSize     = 32 + 32 + 32 + 16 + 16 + 32 + 64 + 8
)

type Block struct {
	Type           byte
	Previous       [32]byte
	Source         [32]byte
	Destination    [32]byte
	Representative [32]byte
	Account        [32]byte
	Link           [32]byte
	Balance        [16]byte
	Amount         [16]byte
	Signature      [64]byte
	Work           [8]byte
}

func (m *Block) ToBlock() blocks.Block {
	common := blocks.CommonBlock{
		Work:      types.Work(m.Work),
		Signature: types.Signature(m.Signature),
	}

	switch m.Type {
	case sendBlock:
		block := blocks.SendBlock{
			types.BlockHash(m.Previous),
			types.AccPub(m.Destination[:]),
			uint128.FromBytes(m.Balance[:]),
			common,
		}
		return &block
	case openBlock:
		block := blocks.OpenBlock{
			types.BlockHash(m.Source),
			types.AccPub(m.Representative[:]),
			types.AccPub(m.Account[:]),
			common,
		}
		return &block
	case changeBlock:
		block := blocks.ChangeBlock{
			types.BlockHash(m.Previous),
			types.AccPub(m.Representative[:]),
			common,
		}
		return &block
	case receiveBlock:
		block := blocks.ReceiveBlock{
			types.BlockHash(m.Previous),
			types.BlockHash(m.Source),
			common,
		}
		return &block
	case utxBlock:
		block := blocks.UtxBlock{
			types.AccPub(m.Account[:]),
			types.BlockHash(m.Previous),
			types.AccPub(m.Representative[:]),
			uint128.FromBytes(m.Balance[:]),
			uint128.FromBytes(m.Amount[:]),
			types.AccPub(m.Link[:]),
			common,
		}
		return &block
	default:
		return nil
	}
}

func (m *Block) Unmarshal(data []byte) error {
	invalidErr := errors.New("invalid block")

	switch m.Type {
	case sendBlock:
		if len(data) != sendSize {
			return invalidErr
		}

		copy(m.Previous[:], data[:32])
		copy(m.Destination[:], data[32:64])
		copy(m.Balance[:], data[64:80])
		copy(m.Signature[:], data[80:144])
		copy(m.Work[:], data[144:152])
	case openBlock:
		if len(data) != openSize {
			return invalidErr
		}

		copy(m.Source[:], data[:32])
		copy(m.Representative[:], data[32:64])
		copy(m.Account[:], data[64:96])
		copy(m.Signature[:], data[96:160])
		copy(m.Work[:], data[160:168])
	case changeBlock:
		if len(data) != changeSize {
			return invalidErr
		}

		copy(m.Previous[:], data[:32])
		copy(m.Representative[:], data[32:64])
		copy(m.Signature[:], data[64:128])
		copy(m.Work[:], data[128:136])
	case receiveBlock:
		if len(data) != receiveSize {
			return invalidErr
		}

		copy(m.Previous[:], data[:32])
		copy(m.Source[:], data[32:64])
		copy(m.Signature[:], data[64:128])
		copy(m.Work[:], data[128:136])
	case utxBlock:
		if len(data) != utxSize {
			return invalidErr
		}

		copy(m.Account[:], data[:32])
		copy(m.Previous[:], data[32:64])
		copy(m.Representative[:], data[64:96])
		copy(m.Balance[:], data[96:112])
		copy(m.Amount[:], data[112:128])
		copy(m.Link[:], data[128:160])
	}

	return nil
}

func (m *Block) Marshal() ([]byte, error) {
	data := make([]byte, 0, 136)

	switch m.Type {
	case sendBlock:
		data = append(data, m.Previous[:]...)
		data = append(data, m.Destination[:]...)
		data = append(data, m.Balance[:]...)
	case openBlock:
		data = append(data, m.Source[:]...)
		data = append(data, m.Representative[:]...)
		data = append(data, m.Account[:]...)
	case changeBlock:
		data = append(data, m.Previous[:]...)
		data = append(data, m.Representative[:]...)
	case receiveBlock:
		data = append(data, m.Previous[:]...)
		data = append(data, m.Source[:]...)
	case utxBlock:
		data = append(data, m.Account[:]...)
		data = append(data, m.Previous[:]...)
		data = append(data, m.Representative[:]...)
		data = append(data, m.Balance[:]...)
		data = append(data, m.Amount[:]...)
		data = append(data, m.Link[:]...)
	}

	data = append(data, m.Signature[:]...)
	data = append(data, m.Work[:]...)

	return data, nil
}
