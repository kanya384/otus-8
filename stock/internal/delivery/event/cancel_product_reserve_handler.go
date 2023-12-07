package event

import (
	"context"

	"stock/internal/domain"
	"stock/pkg/types"
	"stock/pkg/types/products"
)

func (h *Handler) CancelReserveProducts(ctx context.Context, event *products.CancelReserveProducts_v1) (err error) {
	productsList := make([]*domain.ProductReserveItem, 0, len(event.Products))
	for _, pr := range event.Products {
		productsList = append(productsList, &domain.ProductReserveItem{
			ProductId: pr.ProductId,
			Quantity:  pr.Quantity,
		})
	}

	err = h.service.CancelReserveProducts(ctx, event.OrderId, productsList)
	if err != nil {
		return err
	}

	err = h.eventBus.Publish(ctx, products.CancelReserveProductsSuccess_v1{
		Header:  types.NewEventHeader(),
		OrderId: event.OrderId,
	})

	return err
}
