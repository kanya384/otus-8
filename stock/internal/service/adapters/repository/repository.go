package repository

import (
	"context"
	"stock/internal/domain"
)

type Repository interface {
	ReserveProducts(ctx context.Context, products []*domain.ProductReserveItem) error
}
