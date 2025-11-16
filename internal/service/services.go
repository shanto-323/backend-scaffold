package service

import (
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/service/student"
)

type Services struct {
	StudentService student.Service
}

func New(s *server.Server) *Services {
	return &Services{
		StudentService: student.NewService(s),
	}
}
