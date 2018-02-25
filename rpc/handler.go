package rpc

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/frankh/crypto/ed25519"
	"github.com/frankh/nano/address"

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
		"account_get": handlerFn(accountGet),
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

	res["account"] = string(address.PubKeyToAddress(ed25519.PublicKey(pubBytes)))
	json.NewEncoder(w).Encode(res)
}
