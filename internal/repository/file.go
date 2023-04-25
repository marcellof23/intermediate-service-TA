package repository

import (
	"context"

	"github.com/intermediate-service-ta/internal/model"
)

type FileRepository interface {
	Create(c context.Context, file *model.File) (model.File, error)
	GetTotalSizeClient(c context.Context) (map[string]int64, error)
	Delete(c context.Context, filename string) (model.File, error)
	Get(c context.Context, filename string) (model.File, error)
}
