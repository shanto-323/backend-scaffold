package handler

import (
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/service"
)

type Handlers struct {
	services *service.Services
	StudentHandler *Student
}

func NewHandlers(s *server.Server,sr *service.Services) *Handlers {
	return &Handlers{
		StudentHandler: NewStudent(s,sr),
	}
}
