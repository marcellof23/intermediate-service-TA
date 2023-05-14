package publisher

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	watermillMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"

	"github.com/intermediate-service-ta/helper"
	"github.com/intermediate-service-ta/internal/pubsub_notify"
)

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

func GetTraceID(ctx context.Context) string {
	v, ok := ctx.Value("trace_id").(string)
	if !ok {
		return ""
	}

	return v
}

type Publisher struct {
	publisher watermillMessage.Publisher
}

type PublisherInterface interface {
	Publish(ctx context.Context, topic string, message pubsub.Message) error
	init(ctx context.Context, mode pubsub.Mode, loggerMode pubsub.Logger) error
}

var defaultPublisher PublisherInterface

type NoopPublisher struct{}

func (p NoopPublisher) Publish(_ context.Context, _ string, _ pubsub.Message) error {
	return nil
}

func (p NoopPublisher) init(_ context.Context, _ pubsub.Mode, _ pubsub.Logger) error {
	return nil
}

func SetDefault(publisher PublisherInterface) {
	defaultPublisher = publisher
}

func Default() PublisherInterface {
	return defaultPublisher
}

func InitDefault(ctx context.Context) error {
	if defaultPublisher != nil {
		return nil
	}

	mode := pubsub.ModeGoChannel
	if envMode := os.Getenv(pubsub.EnvDefaultMode); envMode != "" {
		mode = pubsub.Mode(envMode)
	}

	logger := pubsub.LoggerNoop
	if envLogger := os.Getenv(pubsub.EnvLogger); envLogger != "" {
		logger = pubsub.Logger(envLogger)
	}

	defaultPublisher = &Publisher{}
	return defaultPublisher.init(ctx, mode, logger)
}

func (p *Publisher) init(ctx context.Context, mode pubsub.Mode, loggerMode pubsub.Logger) error {
	var err error

	var logger watermill.LoggerAdapter
	switch loggerMode {
	case pubsub.LoggerZap:
		logger = &pubsub.ZapLoggerAdapter{
			Logger: helper.GetLogger(ctx),
		}
	case pubsub.LoggerNoop:
		logger = watermill.NopLogger{}
	}

	switch mode {
	case pubsub.ModeGooglePubsub:
		p.publisher, err = googlecloud.NewPublisher(googlecloud.PublisherConfig{
			ProjectID: os.Getenv(pubsub.EnvGoogleProjectID),
		}, logger)
		if err != nil {
			return errors.Wrap(err, "googlecloud.NewPublisher")
		}
	default:
		p.publisher = pubsub.GetOrInitGoChannel(gochannel.Config{}, logger)
	}

	return nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, message pubsub.Message) error {
	topic = os.Getenv(pubsub.EnvPrefix) + "." + topic
	b, err := json.Marshal(message.Payload)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	m := watermillMessage.Message{
		Metadata: map[string]string{
			pubsub.TraceID: GetTraceID(ctx),
		},
		Payload: b,
	}

	if err = p.publisher.Publish(topic, &m); err != nil {
		return errors.Wrap(err, "publisher.Publish")
	}

	return err
}
