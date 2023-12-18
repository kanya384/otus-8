package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewMessageRouter(
	redisPublisher message.Publisher,
	redisSubscriber message.Subscriber,
	eventProcessorConfig cqrs.EventProcessorConfig,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = eventProcessor.AddHandlers()
	if err != nil {
		panic(err)
	}
	return router
}
