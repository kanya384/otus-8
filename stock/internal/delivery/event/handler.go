package event

import (
	"context"
	"stock/internal/domain"

	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(
	service Service,
) Handler {
	if service == nil {
		panic("missing service")
	}

	return Handler{
		service: service,
	}
}

type Service interface {
	ReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error)
	CancelReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) error
}
