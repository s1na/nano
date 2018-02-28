package blocks

import (
	"bytes"
	"encoding/gob"

	"github.com/frankh/nano/store"
	"github.com/frankh/nano/types"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type BlockStore struct {
	s            *store.Store
	orphanBlocks map[types.BlockHash]Block
}

func NewBlockStore(store *store.Store) *BlockStore {
	s := new(BlockStore)

	s.s = store
	s.orphanBlocks = make(map[types.BlockHash]Block)

	// Register block types for gob, so it encodes
	// and decoes the Block interface.
	gob.Register(&OpenBlock{})
	gob.Register(&SendBlock{})
	gob.Register(&ChangeBlock{})
	gob.Register(&ReceiveBlock{})

	return s
}

func (s *BlockStore) Init() error {
	_, err := s.GetBlock(GenesisBlock.Hash())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return s.SetBlock(GenesisBlock)
		} else {
			return err
		}
	}

	return nil
}

func (s *BlockStore) SetBlock(b Block) error {
	if !ValidateBlockWork(b) {
		return errors.New("invalid block work")
	}

	if b.Type() != Open && b.Type() != Change && b.Type() != Send && b.Type() != Receive {
		return errors.New("unknown block type")
	}

	if _, ok := s.orphanBlocks[b.GetPrevious()]; !ok {
		s.orphanBlocks[b.GetPrevious()] = b
		log.WithFields(log.Fields{
			"hash":     b.Hash().String(),
			"previous": b.GetPrevious().String(),
		}).Info("Added orphan block")
		return errors.New("cannot find parent block")
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(b); err != nil {
		return err
	}

	// TODO: Store orphan children
	// TODO: If open, store twice?
	return s.s.Set(append([]byte("block:"), b.Hash().Slice()...), buf.Bytes())
}

func (s *BlockStore) GetBlock(hash types.BlockHash) (Block, error) {
	v, err := s.s.Get(append([]byte("block:"), hash.Slice()...))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	var b Block
	if err = dec.Decode(&b); err != nil {
		return nil, err
	}

	return b, nil
}
