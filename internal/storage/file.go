package storage

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"
	"github.com/intermediate-service-ta/internal/storage/dao"

	_ "github.com/intermediate-service-ta/internal/storage/dao"
)

type filerepository struct{}

func NewFileRepo() repository_intf.FileRepository {
	return &filerepository{}
}

func (fr *filerepository) Create(c context.Context, file *model.File) (model.File, error) {
	var db, err = helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.File{}, err
	}

	if err := db.Create(&file).Error; err != nil {
		return model.File{}, err
	}

	return *file, nil
}

func (fr *filerepository) Get(c *gin.Context, filename string) (model.File, error) {
	var file dao.File
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.File{}, err
	}

	if err := db.Where("filename = ?", filename).First(&file).Error; err != nil {
		return model.File{}, err
	}

	res := dao.ToFileDTO(file)
	return res, nil
}

func (fr *filerepository) Delete(c *gin.Context, filename string) error {
	var file dao.File
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}

	if err := db.Where("filename = ?", filename).Delete(&file).Error; err != nil {
		return err
	}

	return nil
}

func (fr *filerepository) GetTotalSizeClient(c *gin.Context, filename string) error {
	var file dao.File
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}

	if err := db.Where("filename = ?", filename).Delete(&file).Error; err != nil {
		return err
	}

	return nil
}
