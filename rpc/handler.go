package rpc

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/s1na/nano/account"
	"github.com/s1na/nano/blocks"
	"github.com/s1na/nano/types"
	"github.com/s1na/nano/uint128"
	"github.com/s1na/nano/wallet"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type handlerFn func(http.ResponseWriter, *gjson.Result) error

type Handler struct {
	fns map[string]handlerFn
}

func NewHandler() *Handler {
	h := new(Handler)
	h.registerHandlers()
	return h
}

func (h *Handler) registerHandlers() {
	h.fns = map[string]handlerFn{
		"account_get":   handlerFn(accountGet),
		"wallet_create": handlerFn(walletCreate),
		"wallet_add":    handlerFn(walletAdd),
		"send":          handlerFn(send),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respErr(w, "Error in reading request")
		return
	}

	bs := string(b)
	if !gjson.Valid(bs) {
		respErr(w, "Invalid json")
		return
	}

	body := gjson.Parse(bs)
	action := body.Get("action").String()
	handler, ok := h.fns[action]
	if !ok {
		respErr(w, "Action not found")
		return
	}

	if err := handler(w, &body); err != nil {
		respErr(w, err.Error())
	}
}

func accountGet(w http.ResponseWriter, body *gjson.Result) error {
	res := make(map[string]string)

	v := body.Get("key").String()
	pubBytes, err := hex.DecodeString(v)
	if err != nil {
		return errors.New("invalid key")
	}

	res["account"] = types.PubKey(pubBytes).Address()
	json.NewEncoder(w).Encode(res)

	return nil
}

func walletCreate(w http.ResponseWriter, body *gjson.Result) error {
	res := make(map[string]string)

	wal := wallet.NewWallet()
	id, err := wal.Init()
	if err != nil {
		return errors.New("internal error")
	}

	ws := wallet.NewWalletStore(db)
	if err = ws.SetWallet(wal); err != nil {
		return errors.New("internal error")
	}

	walletsCh <- wal

	res["wallet"] = id.Hex()
	json.NewEncoder(w).Encode(res)

	return nil
}

func walletAdd(w http.ResponseWriter, body *gjson.Result) error {
	res := make(map[string]string)

	wid, err := types.PubKeyFromHex(body.Get("wallet").String())
	if err != nil {
		return err
	}

	key, err := types.PrvKeyFromString(body.Get("key").String())
	if err != nil {
		return err
	}

	ws := wallet.NewWalletStore(db)
	wal, err := ws.GetWallet(wid)
	if err != nil {
		return err
	}

	pub, prv, err := types.KeypairFromPrvKey(key)
	if err != nil {
		return err
	}

	wal.InsertAdhoc(pub, prv)
	if err = ws.SetWallet(wal); err != nil {
		return errors.New("internal error")
	}

	walletsCh <- wal

	res["account"] = pub.Address()
	json.NewEncoder(w).Encode(res)

	return nil
}

func send(w http.ResponseWriter, body *gjson.Result) error {
	res := make(map[string]string)

	wid, err := types.PubKeyFromHex(body.Get("wallet").String())
	if err != nil {
		return err
	}

	source, err := types.PubKeyFromAddress(body.Get("source").String())
	if err != nil {
		return err
	}

	dest, err := types.PubKeyFromAddress(body.Get("destination").String())
	if err != nil {
		return err
	}

	auint64, err := strconv.ParseUint(body.Get("amount").String(), 10, 64)
	if err != nil {
		return err
	}
	amount := uint128.FromInts(0, auint64)

	ws := wallet.NewWalletStore(db)
	wal, err := ws.GetWallet(wid)
	if err != nil {
		return errors.New("wallet not found")
	}

	ok := wal.HasAccount(source.Address())
	if !ok {
		return errors.New("account not found")
	}

	as := account.NewAccountStore(db)
	acc, err := as.GetAccount(source)
	if err != nil {
		return errors.New("account info not found")
	}

	b := &blocks.SendBlock{
		Previous:    acc.Head,
		Destination: dest,
		Balance:     acc.Balance.Sub(amount),
		CommonBlock: blocks.CommonBlock{
			Account: source,
		},
	}
	b.Work = types.GenerateWorkForHash(acc.Head)
	b.Signature = wal.Accounts[source.Address()].Sign(b.Hash().Slice())

	blocksCh <- b

	res["sent"] = "true"
	json.NewEncoder(w).Encode(res)

	return nil
}
