package storage

import (
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"
	"github.com/intermediate-service-ta/internal/storage/dao"
)

type userrepository struct{}

func NewUserRepo() repository_intf.UserRepository {
	return &userrepository{}
}

func (ur *userrepository) Create(c *gin.Context, user *model.User) (model.User, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.User{}, err
	}

	if err := db.Create(&user).Error; err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (ur *userrepository) FindByUsername(c *gin.Context, username string) (model.User, error) {
	var user dao.User
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.User{}, err
	}

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return model.User{}, err
	}

	res := dao.ToUserDTO(user)
	return res, nil
}

func (ur *userrepository) Update(c *gin.Context, user *model.User) (model.User, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.User{}, err
	}

	if err := db.Updates(&user).Error; err != nil {
		return model.User{}, err
	}

	return *user, nil
}
