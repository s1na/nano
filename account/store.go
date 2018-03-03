package account

import (
	"bytes"
	"encoding/gob"

	"github.com/frankh/nano/store"
	"github.com/frankh/nano/types"
)

type AccountStore struct {
	s *store.Store
}

func NewAccountStore(store *store.Store) *AccountStore {
	s := new(AccountStore)

	s.s = store

	return s
}

func (s *AccountStore) SetAccount(a *Account) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(a); err != nil {
		return err
	}

	return s.s.Set(append([]byte("account:"), a.PublicKey...), buf.Bytes())
}

func (s *AccountStore) GetAccount(pub types.PubKey) (*Account, error) {
	v, err := s.s.Get(append([]byte("account:"), pub...))
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(v)
	dec := gob.NewDecoder(buf)

	var a *Account
	if err = dec.Decode(&a); err != nil {
		return nil, err
	}

	return a, nil
}
