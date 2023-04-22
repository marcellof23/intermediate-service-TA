package repository

import (
	"context"

	"github.com/intermediate-service-ta/internal/model"
)

type FileRepository interface {
	Create(c context.Context, file *model.File) (model.File, error)
}
