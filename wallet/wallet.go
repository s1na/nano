package wallet

import (
	"encoding/hex"
	"strings"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/account"
)

type Wallet struct {
	Seed     string
	Accounts map[string]*account.Account
}

func NewWallet() *Wallet {
	w := new(Wallet)

	w.Accounts = make(map[string]*account.Account)

	return w
}

func (w *Wallet) GenerateSeed() (string, error) {
	pub, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", err
	}

	w.Seed = strings.ToUpper(hex.EncodeToString(pub))

	return w.Seed, nil
}
