package service

import "github.com/shanto-323/backend-scaffold/internal/service/student"

type Services struct {
	StudentService student.Service
}

func New() *Services {
	return &Services{
		StudentService: student.NewService(),
	}
}
