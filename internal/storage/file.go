package storage

import (
	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"

	_ "github.com/intermediate-service-ta/internal/storage/dao"
)

type filerepository struct{}

func NewFileRepo() repository_intf.FileRepository {
	return &filerepository{}
}

func (fr *filerepository) Create(c *gin.Context, file *model.File) (model.File, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.File{}, err
	}

	file.Url = "/" + "/" + file.Filename

	if err := db.Create(&file).Error; err != nil {
		return model.File{}, err
	}

	return *file, nil
}
