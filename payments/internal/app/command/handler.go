package command

import (
	"context"
	"payments/internal/domain/payment"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type Handler struct {
	paymentsRepository PaymentsRepository
}

func NewHandler(
	eventBus *cqrs.EventBus,
	paymentsRepository PaymentsRepository,
) Handler {
	if eventBus == nil {
		panic("eventBus is required")
	}
	if paymentsRepository == nil {
		panic("payments repository is required")
	}

	handler := Handler{
		paymentsRepository: paymentsRepository,
	}

	return handler
}

type PaymentsRepository interface {
	AddPayment(ctx context.Context, pmnt *payment.Payment) (*payment.Payment, error)
	ReadPayment(ctx context.Context, paymentUUID string) (*payment.Payment, error)
}
