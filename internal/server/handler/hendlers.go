package handler

import (
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/service"
)

type Handlers struct {
	HealthHandler  *HealthHandler
	StudentHandler *Student
}

func New(s *server.Server, sr *service.Services) *Handlers {
	return &Handlers{
		HealthHandler:  NewHealthHandler(s),
		StudentHandler: NewStudent(s, sr),
	}
}
