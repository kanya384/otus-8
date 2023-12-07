package event

import (
	"context"
	"stock/internal/domain"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type Handler struct {
	service  Service
	eventBus *cqrs.EventBus
}

func NewHandler(
	service Service,
	eventBus *cqrs.EventBus,
) Handler {
	if service == nil {
		panic("missing service")
	}

	return Handler{
		service:  service,
		eventBus: eventBus,
	}
}

type Service interface {
	ReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error)
	CancelReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) error
}
