package rpc

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	s       *http.Server
	handler *Handler
}

func NewServer() *Server {
	s := new(Server)

	s.handler = NewHandler()
	s.s = &http.Server{
		Addr:    ":7076",
		Handler: s.handler,
	}

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
