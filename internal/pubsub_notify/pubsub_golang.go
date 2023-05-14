package pubsub

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"sync"
)

var defaultGoChannel *gochannel.GoChannel
var defaultGoChannelOnce sync.Once

func GetOrInitGoChannel(config gochannel.Config, logger watermill.LoggerAdapter) *gochannel.GoChannel {
	defaultGoChannelOnce.Do(func() {
		defaultGoChannel = gochannel.NewGoChannel(config, logger)
	})

	return defaultGoChannel
}
