package service

import (
	"context"
	"stock/internal/domain"
	"stock/internal/service/adapters/repository"
	"stock/pkg/types"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
)

type service struct {
	repository repository.Repository
	eventBus   *cqrs.EventBus
}

func (s *service) ReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error) {
	err = s.repository.ReserveProducts(ctx, products)
	if err != nil {
		err = s.eventBus.Publish(ctx, domain.ReserveProductsFail_v1{
			Header:  types.NewEventHeader(),
			OrderId: orderId,
			Reason:  err.Error(),
		})
		return err
	}

	err = s.eventBus.Publish(ctx, domain.ReserveProductsSucces_v1{
		Header:  types.NewEventHeader(),
		OrderId: orderId,
	})
	return err
}

func (s *service) CancelReserveProducts(ctx context.Context, orderId uuid.UUID, products []*domain.ProductReserveItem) (err error) {
	err = s.repository.ReserveProducts(ctx, products)
	if err != nil {
		return err
	}

	err = s.eventBus.Publish(ctx, domain.CancelReserveProductsSuccess_v1{
		Header:  types.NewEventHeader(),
		OrderId: orderId,
	})
	return err
}
