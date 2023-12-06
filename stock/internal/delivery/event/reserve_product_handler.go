package event

import (
	"context"
	"stock/internal/domain"
)

func (h *Handler) ReserveProducts(ctx context.Context, event *domain.ReserveProducts_v1) (err error) {
	products := make([]*domain.ProductReserveItem, 0, len(event.Products))
	for _, pr := range event.Products {
		products = append(products, &domain.ProductReserveItem{
			ProductId: pr.ProductId,
			Quantity:  pr.Quantity,
		})
	}

	err = h.service.ReserveProducts(ctx, event.OrderId, products)

	return err
}
