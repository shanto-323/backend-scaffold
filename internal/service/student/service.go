package student

import (
	"context"

	"github.com/shanto-323/backend-scaffold/model"
)

type Service interface {
	Create(ctx context.Context, payload *model.Student) (*model.Student, error)
}
