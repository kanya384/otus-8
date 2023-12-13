package command

import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	eventBus         *cqrs.EventBus
	ordersRepository OrdersRepository
}

func NewHandler(
	eventBus *cqrs.EventBus,
	ordersRepository OrdersRepository,
) Handler {
	if eventBus == nil {
		panic("eventBus is required")
	}
	if ordersRepository == nil {
		panic("orders repository is required")
	}

	handler := Handler{
		eventBus:         eventBus,
		ordersRepository: ordersRepository,
	}

	return handler
}

type OrdersRepository interface {
}
