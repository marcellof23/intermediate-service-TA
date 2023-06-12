package cron

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"google.golang.org/api/option"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/repository"
)

type Cron struct {
	ctx            context.Context
	subscriberRepo repository.SubscriberRepository
}

func NewCronHandler(subscriberRepo repository.SubscriberRepository) *Cron {
	return &Cron{subscriberRepo: subscriberRepo}
}

func GetTopic(ctx context.Context, c *pubsub.Client, topic string) *pubsub.Topic {
	log, ok := ctx.Value("server-logger").(*log.Logger)
	if !ok {
		return &pubsub.Topic{}
	}

	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		return t
	}
	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		log.Fatalf("Failed to create the topic: %v", err)
	}
	return t
}

func (cr *Cron) CronJob() {
	totalUsed, err := cr.subscriberRepo.CountInUsed(cr.ctx)
	if err != nil {
		fmt.Print(err)
	}

	var fraction float64
	fraction = 0
	if totalUsed.Total != 0 {
		fraction = float64(totalUsed.NumberInused / totalUsed.Total)
	}

	if fraction > 0.8 || (fraction == 0 && totalUsed.Total == 0) {
		conf, err := helper.GetConfigFromContext(cr.ctx)
		if err != nil {
			fmt.Print(err)
		}

		proj := conf.Pubsub.Project
		topic := conf.Pubsub.Topic
		credentialsFile := conf.Pubsub.CredentialFile

		client, err := pubsub.NewClient(cr.ctx, proj, option.WithCredentialsFile(credentialsFile))
		if err != nil {
			log.Fatalf("Could not create pubsub Client: %v", err)
		}

		t := GetTopic(cr.ctx, client, topic)

		for i := 0; i < 20; i++ {
			subID := "command-" + uuid.New().String()
			sub, err := client.CreateSubscription(cr.ctx, subID, pubsub.SubscriptionConfig{
				ExpirationPolicy:      120 * time.Hour,
				Topic:                 t,
				AckDeadline:           20 * time.Second,
				EnableMessageOrdering: true,
			})
			if err != nil {
				fmt.Printf("Could not create pubsub Client: %v", err)
			}

			cr.subscriberRepo.Create(cr.ctx, sub.ID())
		}

	}
}

func (cr *Cron) ManageSubscriber(ctx context.Context) error {
	c := cron.New()
	cr.ctx = ctx
	err := c.AddFunc("@every 1m", cr.CronJob)
	if err != nil {
		return err
	}

	c.Start()
	return nil
}
