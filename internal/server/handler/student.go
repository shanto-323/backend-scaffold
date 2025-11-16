package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/internal/service"
	"github.com/shanto-323/backend-scaffold/model"
)

type Student struct {
	s  *server.Server
	sr *service.Services
}

func NewStudent(s *server.Server, sr *service.Services) *Student {
	return &Student{
		s:  s,
		sr: sr,
	}
}

func (stud *Student) Create(c echo.Context) error {
	return Handle(
		func(c echo.Context, payload *model.Student) (*model.Student, error) {
			return stud.sr.StudentService.Create(c.Request().Context(), payload)
		},
		http.StatusCreated,
		&model.Student{},
	)(c)
}
