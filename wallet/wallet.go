package wallet

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/account"
	"github.com/frankh/nano/address"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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

func (w *Wallet) NewAccount() *account.Account {
	a := account.NewAccount()

	pub, prv := address.KeypairFromSeed(w.Seed, uint32(len(w.Accounts)))
	a.PublicKey = pub
	a.PrivateKey = prv
	w.Accounts[string(a.Address())] = a

	return a
}

func (w *Wallet) String() string {
	b, err := json.Marshal(w)
	if err != nil {
		log.Warn(err)
		return w.Id
	}

	return string(b)
}
