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
	AddPayment(ctx context.Context, payment *Payment) (*Payment, error)

	ReadPayment(ctx context.Context, paymentUUID string) (*Payment, error)

	UpdatePayment(
		ctx context.Context,
		paymentUUID string,
		updateFn func(ctx context.Context, oldPayemnt *Payment) (*Payment, error),
	) (*Payment, error)
}
