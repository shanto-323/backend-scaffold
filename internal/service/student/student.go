package student

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shanto-323/backend-scaffold/internal/server"
	"github.com/shanto-323/backend-scaffold/model"
	"go.opentelemetry.io/otel/attribute"
)

type student struct {
	s *server.Server
}

func NewService(s *server.Server) Service {
	return &student{
		s: s,
	}
}

func (st *student) Create(ctx context.Context, payload *model.Student) (*model.Student, error) {
	_, span := st.s.TraceProvider.Tracer.Start(ctx, "service")
	defer span.End()

	start := time.Now()
	defer func() {
		totalTime := time.Since(start)
		span.SetAttributes(
			attribute.String("total", totalTime.String()),
		)

	}()

	time.Sleep(2 * time.Second)

	return &model.Student{
		ID:   uuid.New(),
		Name: payload.Name,
		Roll: payload.Roll,
	}, nil
}
