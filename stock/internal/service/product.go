package service

import (
	"context"
	"stock/internal/domain"
	"stock/internal/service/adapters/repository"

	"github.com/google/uuid"
)

type Service struct {
	repository repository.Repository
}

func New(
	repository repository.Repository,
) Service {
	return Service{
		repository: repository,
	}
}

func (s *Service) CreateProduct(ctx context.Context, product domain.Product) (err error) {
	return s.repository.CreateProduct(ctx, product)
}

func (s *Service) ReadProducts(ctx context.Context) (products []domain.Product, err error) {
	return s.repository.AllProduct(ctx)
}

func (s *Service) ReadProductById(ctx context.Context, productId uuid.UUID) (product domain.Product, err error) {
	return s.repository.ProductByID(ctx, productId)
}

func (s *Service) ReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error) {
	err = s.repository.ReserveProducts(ctx, products)
	return err
}

func (s *Service) CancelReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error) {
	err = s.repository.ReserveProducts(ctx, products)
	return err
}
