package wallet

import (
	"bytes"
	"encoding/gob"

	"github.com/frankh/nano/store"
)

type WalletStore struct {
	s *store.Store
}

func NewWalletStore(store *store.Store) *WalletStore {
	s := new(WalletStore)

	s.s = store

	return s
}

func (s *WalletStore) SetWallet(w *Wallet) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(w); err != nil {
		return err
	}

	return s.s.Set([]byte("wallet:"+w.Id), buf.Bytes())
}

func (s *WalletStore) GetWallet(id string) (*Wallet, error) {
	v, err := s.s.Get([]byte("wallet:" + id))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	var w *Wallet
	if err = dec.Decode(&w); err != nil {
		return nil, err
	}

	return w, nil
}
