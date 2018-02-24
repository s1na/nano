package wallet

import (
	"bytes"
	"encoding/gob"

	"github.com/frankh/nano/account"
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

func (s *WalletStore) GetWallets() (map[string]*Wallet, error) {
	res, err := s.s.GetPrefixValues([]byte("wallet:"))
	if err != nil {
		return nil, err
	}

	wallets := make(map[string]*Wallet)
	for k, v := range res {
		buf := bytes.NewBuffer(v)
		dec := gob.NewDecoder(buf)

		var w *Wallet
		if err = dec.Decode(&w); err != nil {
			return nil, err
		}

		wallets[string(k)] = w
	}

	return wallets, nil
}

func (s *WalletStore) GetAccount(addr string) (*account.Account, error) {
	wallets, err := s.GetWallets()
	if err != nil {
		return nil, err
	}

	for _, v := range wallets {
		if a, exists := v.Accounts[addr]; exists {
			return a, nil
		}
	}

	return nil, nil
}
