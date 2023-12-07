package event

import (
	"context"

	"stock/internal/domain"
	"stock/pkg/types"
	"stock/pkg/types/products"
)

func (h *Handler) ReserveProducts(ctx context.Context, event *products.ReserveProducts_v1) (err error) {
	productsList := make([]*domain.ProductReserveItem, 0, len(event.Products))
	for _, pr := range event.Products {
		productsList = append(productsList, &domain.ProductReserveItem{
			ProductId: pr.ProductId,
			Quantity:  pr.Quantity,
		})
	}

	err = h.service.ReserveProducts(ctx, event.OrderId, productsList)
	if err != nil {
		err = h.eventBus.Publish(ctx, products.ReserveProductsFail_v1{
			Header:  types.NewEventHeader(),
			OrderId: event.OrderId,
			Reason:  err.Error(),
		})
		return err
	}

	err = h.eventBus.Publish(ctx, products.ReserveProductsSucces_v1{
		Header:  types.NewEventHeader(),
		OrderId: event.OrderId,
	})

	return err
}
