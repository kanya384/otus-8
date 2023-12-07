package event

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewRouter(
	eventHandler Handler,
	eventProcessorConfig cqrs.EventProcessorConfig,
	watermillLogger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		panic(err)
	}

	useMiddlewares(router, watermillLogger)

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(router, eventProcessorConfig)
	if err != nil {
		panic(err)
	}

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"ReserveProducts",
			eventHandler.ReserveProducts,
		),
		cqrs.NewEventHandler(
			"CancelReserveProducts",
			eventHandler.CancelReserveProducts,
		),
	)

	if err != nil {
		panic(err)
	}

	return router
}
