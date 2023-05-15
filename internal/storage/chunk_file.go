package storage

import (
	"context"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"
	"github.com/intermediate-service-ta/internal/storage/dao"
)

type chunkfilerepository struct{}

func NewChunkFileRepo() repository_intf.ChunkFileRepository {
	return &chunkfilerepository{}
}

func (fcr *chunkfilerepository) Create(c context.Context, chunkFile *model.ChunkFile) (model.ChunkFile, error) {
	var db, err = helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.ChunkFile{}, err
	}

	if err := db.Create(&chunkFile).Error; err != nil {
		return model.ChunkFile{}, err
	}

	return *chunkFile, nil
}

func (fcr *chunkfilerepository) Get(c context.Context, filename string) (model.ChunkFile, error) {
	var chunkFile dao.ChunkFile
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.ChunkFile{}, err
	}

	if err := db.Where("filename = ?", filename).Find(&chunkFile).Order("order asc").Error; err != nil {
		return model.ChunkFile{}, err
	}

	res := dao.ToFileChunkDTO(chunkFile)
	return res, nil
}
