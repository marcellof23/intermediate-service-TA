package storage

import (
	"context"
	"sync"

	"github.com/intermediate-service-ta/boot"
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

var TotalSizeClient = make(map[string]int64, 0)

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

func (fr *filerepository) Get(c context.Context, filename string) (model.File, error) {
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

func (fr *filerepository) Delete(c context.Context, filename string) (model.File, error) {
	var file, getFile dao.File
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.File{}, err
	}

	if err := db.Where("filename = ?", filename).First(&getFile).Error; err != nil {
		return model.File{}, err
	}

	if err := db.Where("filename = ?", filename).Delete(&file).Error; err != nil {
		return model.File{}, err
	}

	return dao.ToFileDTO(getFile), nil
}

type ResultClientSize struct {
	Client    string
	TotalSize int64
}

func (fr *filerepository) GetTotalSizeClient(c context.Context) (map[string]int64, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return map[string]int64{}, err
	}

	var res []ResultClientSize
	db.Raw("SELECT client, sum(size) as total_size FROM files group by client").Scan(&res)

	m := make(map[string]int64, 0)
	for _, v := range res {
		m[v.Client] = v.TotalSize
	}

	for _, val := range boot.Clients {
		_, isKeyPresent := m[val]
		if !isKeyPresent {
			m[val] = 0
		}
	}

	return m, nil
}

func UpdateTotalSizeClient(client string, size int64) {
	var mtx sync.Mutex
	mtx.Lock()
	TotalSizeClient[client] += size
	mtx.Unlock()
}
