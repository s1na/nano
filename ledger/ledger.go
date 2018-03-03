package ledger

import (
	"github.com/frankh/nano/account"
	"github.com/frankh/nano/blocks"
	"github.com/frankh/nano/store"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

type Ledger struct {
	store *store.Store
	bs    *blocks.BlockStore
	as    *account.AccountStore
}

func NewLedger(s *store.Store) *Ledger {
	l := new(Ledger)

	l.store = s
	l.bs = blocks.NewBlockStore(s)
	l.as = account.NewAccountStore(s)

	return l
}

func (l *Ledger) Init() error {
	_, err := l.bs.GetBlock(blocks.GenesisBlock.Hash())
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return l.AddOpen(blocks.GenesisBlock)
		} else {
			return err
		}
	}

	return nil

}

func (l *Ledger) AddSend(b *blocks.SendBlock) error {
	if err := l.bs.SetBlock(b); err != nil {
		return err
	}

	acc, err := l.as.GetAccount(b.Account)
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil
		}

		return err
	}

	acc.Balance = b.Balance
	acc.Head = b.Hash()
	if err = l.as.SetAccount(acc); err != nil {
		return err
	}

	log.Printf("Added block %s for account %s\n", b.Hash(), acc.Address())

	return nil
}

func (l *Ledger) AddOpen(b *blocks.OpenBlock) error {
	if err := l.bs.SetBlock(b); err != nil {
		return err
	}

	acc := account.NewAccount()
	acc.PublicKey = b.Account
	acc.Rep = b.Representative
	acc.Head = b.Hash()
	acc.Open = b.Hash()

	if b.Hash() == blocks.GenesisBlock.Hash() {
		acc.Balance = blocks.GenesisAmount
	}

	if err := l.as.SetAccount(acc); err != nil {
		return err
	}

	log.Printf("Added block %s for account %s\n", b.Hash(), acc.Address())

	return nil
}

func (l *Ledger) AddBlock(block blocks.Block) error {
	switch b := block.(type) {
	case *blocks.SendBlock:
		return l.AddSend(b)
	case *blocks.OpenBlock:
		return l.AddOpen(b)
	}

	return nil
}
