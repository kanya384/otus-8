package message

import (
	"order/internal/app/command/saga"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func NewMessageRouter(
	redisPublisher message.Publisher,
	redisSubscriber message.Subscriber,
	saga *saga.OrderSaga,
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

	err = eventProcessor.AddHandlers(
		cqrs.NewEventHandler(
			"SagaOnOrderInitialized",
			saga.OnOrderInitialized,
		),
		cqrs.NewEventHandler(
			"SagaOnPaymentConfirmed",
			saga.OnPaymentConfirmed,
		),
		cqrs.NewEventHandler(
			"SagaOnPaymentFailed",
			saga.OnPaymentFailed,
		),
		cqrs.NewEventHandler(
			"SagaOnOrderProductsReserved",
			saga.OnOrderProductsReserved,
		),
		cqrs.NewEventHandler(
			"SagaOnOrderProductsReserveFailed",
			saga.OnOrderProductsReserveFailed,
		),
	)
	if err != nil {
		panic(err)
	}
	return router
}
