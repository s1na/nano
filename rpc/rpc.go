package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/s1na/nano/blocks"
	"github.com/s1na/nano/store"
	"github.com/s1na/nano/wallet"

	log "github.com/sirupsen/logrus"
)

var (
	db        *store.Store
	walletsCh chan *wallet.Wallet
	blocksCh  chan blocks.Block
)

type Server struct {
	s       *http.Server
	handler *Handler
}

func NewServer(st *store.Store, wCh chan *wallet.Wallet, bCh chan blocks.Block) *Server {
	s := new(Server)

	s.handler = NewHandler()
	s.s = &http.Server{
		Addr:    ":7076",
		Handler: s.handler,
	}
	db = st
	walletsCh = wCh
	blocksCh = bCh

	return s
}

func (s *Server) Start() {
	go func() {
		if err := s.s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) Stop() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	s.s.Shutdown(ctx)
	log.Info("RPC server gracefully stopped")
}

func respErr(w http.ResponseWriter, msg string) {
	res := map[string]string{"error": msg}
	json.NewEncoder(w).Encode(res)
}
