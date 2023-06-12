package storage

import (
	"context"
	"sync"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"
	"github.com/intermediate-service-ta/internal/storage/dao"
)

type requestrepository struct {
	Mu sync.Mutex
}

func NewRequestRepo() repository_intf.RequestRepository {
	return &requestrepository{}
}

func (rr *requestrepository) Create(c context.Context, reqcommand *model.Request) (model.Request, error) {
	var db, err = helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.Request{}, err
	}

	if err := db.Create(&reqcommand).Error; err != nil {
		return model.Request{}, err
	}

	return *reqcommand, nil
}

type ResultTotalCommand struct {
	TotalCommand         int
	TotalSuccessExecuted int
}

func (rr *requestrepository) IsSuccess(c context.Context, reqID string) (bool, error) {
	var req dao.Request
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return false, err
	}

	rr.Mu.Lock()
	if err := db.Where("request_id = ?", reqID).First(&req).Error; err != nil {
		return false, err
	}
	rr.Mu.Unlock()

	if req.TotalCommand == req.TotalSuccessExecuted {
		return true, nil
	}

	return false, nil
}

func (rr *requestrepository) AddExecutedCommand(c context.Context, reqID string) error {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}
	var req dao.Request

	rr.Mu.Lock()
	if err := db.Where("request_id = ?", reqID).First(&req).Error; err != nil {
		return err
	}

	totalExec := req.TotalSuccessExecuted
	totalExec += 1

	err = db.Exec("UPDATE request SET total_success_executed = ? WHERE request_id = ?", totalExec, reqID).Error
	if err != nil {
		return err
	}
	rr.Mu.Unlock()

	return nil
}
