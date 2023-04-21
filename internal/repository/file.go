package repository

import (
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/internal/model"
)

type FileRepository interface {
	Create(c *gin.Context, file *model.File) (model.File, error)
}
