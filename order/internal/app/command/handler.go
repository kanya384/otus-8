package command

import (
	"context"
	"order/internal/domain/order"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
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
		ordersRepository: ordersRepository,
	}

	return handler
}

type OrdersRepository interface {
	AddOrder(ctx context.Context, or *order.Order) (*order.Order, error)
	ReadOrder(ctx context.Context, orderUUID string) (*order.Order, error)
}
