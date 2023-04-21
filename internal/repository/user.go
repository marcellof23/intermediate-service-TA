package repository

import (
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/internal/model"
)

type UserRepository interface {
	Create(c *gin.Context, user *model.User) (model.User, error)
	FindByUsername(c *gin.Context, username string) (model.User, error)
}
