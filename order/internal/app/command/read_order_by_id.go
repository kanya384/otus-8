package command

import (
	"context"
	"order/internal/domain/order"
)

type ReadOrder struct {
	OrderUUID string
}

func (h Handler) ReadOrderById(ctx context.Context, cmd ReadOrder) (order *order.Order, err error) {

	order, err = h.ordersRepository.ReadOrder(ctx, cmd.OrderUUID)
	if err != nil {
		return
	}
	return
}
