package wallet

import (
	"encoding/json"

	"github.com/s1na/nano/types"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Wallet struct {
	Id       types.PubKey
	Seed     types.PrvKey
	Accounts map[string]types.PrvKey
}

func NewWallet() *Wallet {
	w := new(Wallet)

	w.Accounts = make(map[string]types.PrvKey)

	return w
}

func (w *Wallet) Init() (types.PubKey, error) {
	if err := w.GenerateSeed(); err != nil {
		return types.PubKey{}, errors.Wrap(err, "generating seed failed")
	}

	if err := w.GenerateID(); err != nil {
		return types.PubKey{}, errors.Wrap(err, "generating id failed")
	}

	return w.Id, nil
}

func (w *Wallet) GenerateSeed() error {
	_, prv, err := types.GenerateKey(nil)
	if err != nil {
		return err
	}

	w.Seed = prv

	return nil
}

func (w *Wallet) GenerateID() error {
	pub, _, err := types.GenerateKey(nil)
	if err != nil {
		return err
	}

	w.Id = pub

	return nil
}

func (w *Wallet) NewAccount() types.PubKey {
	pub, prv, _ := types.KeypairFromSeed(w.Seed, uint32(len(w.Accounts)))
	w.Accounts[pub.Address()] = prv

	return pub
}

func (w *Wallet) InsertAdhoc(pub types.PubKey, prv types.PrvKey) {
	w.Accounts[pub.Address()] = prv
}

func (w *Wallet) HasAccount(addr string) bool {
	_, ok := w.Accounts[addr]
	return ok
}

/*func (w *Wallet) GetAccount(addr string) (*account.Account, bool) {
	a, ok := w.Accounts[addr]
	if !ok {
		return nil, false
	}

	return a, true
}*/

func (w *Wallet) String() string {
	b, err := json.Marshal(w)
	if err != nil {
		log.Warn(err)
		return w.Id.Hex()
	}

	return string(b)
}
