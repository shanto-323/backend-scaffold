package handler

import "github.com/shanto-323/backend-scaffold/internal/server"

type Handlers struct{}

func NewHandlers(s *server.Server) *Handlers {
	return &Handlers{}
}

