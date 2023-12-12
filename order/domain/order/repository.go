package order

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	OrderUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order '%s' not found", e.OrderUUID)
}

type Repository interface {
	AddOrder(ctx context.Context, or *Order) error

	ReadOrder(ctx context.Context, orderUUID string) (*Order, error)

	UpdateOrder(
		ctx context.Context,
		orderUUID string,
		updateFn func(ctx context.Context, or *Order) (*Order, error),
	) error
}
