package storage

import (
	"context"
	"errors"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/model"
	repository_intf "github.com/intermediate-service-ta/internal/repository"
	"github.com/intermediate-service-ta/internal/storage/dao"
)

type subscriberrepository struct {
}

func NewSubscriberRepo() repository_intf.SubscriberRepository {
	return &subscriberrepository{}
}

func (sr *subscriberrepository) CreateBulk(c context.Context, subs []string) ([]model.Subscriber, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return []model.Subscriber{}, err
	}

	var arrSubs []dao.Subscriber
	for _, val := range subs {
		arrSubs = append(arrSubs, dao.Subscriber{
			GoogleSubscriberID: val,
		})
	}

	if err := db.CreateInBatches(arrSubs, 100).Error; err != nil {
		return []model.Subscriber{}, err
	}

	var modelSubs []model.Subscriber
	for _, val := range arrSubs {
		modelSubs = append(modelSubs, model.Subscriber{
			ID:                 val.ID,
			GoogleSubscriberID: val.GoogleSubscriberID,
			InUsed:             val.InUsed,
			Base: model.Base{
				CreatedAt: val.CreatedAt,
				UpdatedAt: val.UpdatedAt,
			},
		})
	}

	return modelSubs, nil
}

func (sr *subscriberrepository) Create(c context.Context, subID string) error {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}

	subs := dao.Subscriber{
		GoogleSubscriberID: subID,
		InUsed:             false,
	}

	if err := db.Create(&subs).Error; err != nil {
		return err
	}

	return nil
}

func (sr *subscriberrepository) GetSubscription(c context.Context) (string, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return "", err
	}

	var subs dao.Subscriber
	if err := db.Where("in_used = 0").First(&subs).Error; err != nil {
		return "", errors.New("no subscriber available")
	}

	return subs.GoogleSubscriberID, nil
}

func (sr *subscriberrepository) CountInUsed(c context.Context) (model.TotalInUsed, error) {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return model.TotalInUsed{}, err
	}

	var res model.TotalInUsed
	err = db.Raw("SELECT count(in_used) as total FROM subscribers where in_used=1").Scan(&res.NumberInused).Error
	if err != nil {
		return model.TotalInUsed{}, err
	}

	err = db.Raw("SELECT count(*) as total FROM subscribers").Scan(&res.Total).Error
	if err != nil {
		return model.TotalInUsed{}, err
	}

	return res, nil
}

func (sr *subscriberrepository) UpdateSubscriber(c context.Context, subsID string, val bool) error {
	db, err := helper.GetDatabaseFromContext(c) // Get model if exist
	if err != nil {
		return err
	}

	var subs dao.Subscriber
	if err := db.Model(&subs).Where("google_subscriber_id = ?", subsID).Update("in_used", val).Error; err != nil {
		return err
	}

	return nil
}
