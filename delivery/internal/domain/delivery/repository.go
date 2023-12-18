package payment

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
	AddDelivery(ctx context.Context, payment *Delivery) (*Delivery, error)

	ReadDelivery(ctx context.Context, paymentUUID string) (*Delivery, error)

	UpdateDelivery(
		ctx context.Context,
		paymentUUID string,
		updateFn func(ctx context.Context, oldPayemnt *Delivery) (*Delivery, error),
	) (*Delivery, error)
}
