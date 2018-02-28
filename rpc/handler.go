package rpc

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/frankh/nano/types"
	"github.com/frankh/nano/wallet"

	"github.com/tidwall/gjson"
)

type handlerFn func(http.ResponseWriter, *gjson.Result)

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
	action := body.Get("action_type").String()
	handler, ok := h.fns[action]
	if !ok {
		respErr(w, "Action not found")
		return
	}

	handler(w, &body)
}

func accountGet(w http.ResponseWriter, body *gjson.Result) {
	res := make(map[string]string)

	v := body.Get("key").String()
	pubBytes, err := hex.DecodeString(v)
	if err != nil {
		res["error"] = "Invalid key"
		json.NewEncoder(w).Encode(res)
		return
	}

	res["account"] = types.AccPub(pubBytes).String()
	json.NewEncoder(w).Encode(res)
}

func walletCreate(w http.ResponseWriter, body *gjson.Result) {
	res := make(map[string]string)

	wal := wallet.NewWallet()
	id, err := wal.Init()
	if err != nil {
		res["error"] = "Internal error"
		json.NewEncoder(w).Encode(res)
		return
	}

	ws := wallet.NewWalletStore(db)
	if err = ws.SetWallet(wal); err != nil {
		res["error"] = "Internal error"
		json.NewEncoder(w).Encode(res)
		return
	}

	walletsCh <- wal

	res["wallet"] = id
	json.NewEncoder(w).Encode(res)
}
