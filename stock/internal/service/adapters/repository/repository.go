package repository

import (
	"context"
	"stock/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	CreateProduct(ctx context.Context, product domain.Product) error
	ReserveProducts(ctx context.Context, products []*domain.ProductReserveItem) error
	AllProduct(ctx context.Context) ([]domain.Product, error)
	ProductByID(ctx context.Context, productID uuid.UUID) (domain.Product, error)
}
