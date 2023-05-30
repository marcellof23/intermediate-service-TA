package repository

import (
	"context"

	"github.com/intermediate-service-ta/internal/model"
)

type SubscriberRepository interface {
	CreateBulk(c context.Context, subs []string) ([]model.Subscriber, error)
	Create(c context.Context, subID string) error
	GetSubscription(c context.Context) (string, error)
	CountInUsed(c context.Context) (model.TotalInUsed, error)
	UpdateSubscriber(c context.Context, subsID string, val bool) error
}
