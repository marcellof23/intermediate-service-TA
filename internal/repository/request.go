package repository

import (
	"context"

	"github.com/intermediate-service-ta/internal/model"
)

type RequestRepository interface {
	Create(c context.Context, reqcommand *model.Request) (model.Request, error)
}
