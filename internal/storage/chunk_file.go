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

	var file dao.File
	if err := db.Where("id = ?", chunkFile.FileID).First(&file).Error; err != nil {
		return model.ChunkFile{}, err
	}

	file.Size = file.Size + chunkFile.Size
	if err := db.Where("id = ?", chunkFile.FileID).Updates(&file).Error; err != nil {
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

func (fcr *chunkfilerepository) GetChunkFileByFileID(c context.Context, fid int64) ([]model.ChunkFile, error) {
	var chunkFiles []dao.ChunkFile
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return []model.ChunkFile{}, err
	}

	if err := db.Where("file_id = ?", fid).Find(&chunkFiles).Error; err != nil {
		return []model.ChunkFile{}, err
	}

	var res []model.ChunkFile
	for _, cf := range chunkFiles {
		res = append(res, dao.ToFileChunkDTO(cf))
	}

	return res, nil
}

func (fcr *chunkfilerepository) DeleteChunkFileByFileID(c context.Context, fid int64) error {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}

	if err := db.Where("file_id = ?", fid).Delete(&dao.ChunkFile{}).Error; err != nil {
		return err
	}

	return nil
}
