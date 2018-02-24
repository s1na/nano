package wallet

import (
	"encoding/hex"
	"strings"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/account"

	"github.com/pkg/errors"
)

type Wallet struct {
	Id       string
	Seed     string
	Accounts map[string]*account.Account
}

func NewWallet() *Wallet {
	w := new(Wallet)

	w.Accounts = make(map[string]*account.Account)

	return w
}

func (w *Wallet) Init() (string, error) {
	_, prv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", errors.Wrap(err, "generating seed failed")
	}

	w.Seed = strings.ToUpper(hex.EncodeToString(prv))

	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", errors.Wrap(err, "generating id failed")
	}

	w.Id = strings.ToUpper(hex.EncodeToString(pub))

	return w.Id, nil
}
