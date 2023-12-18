package command

import (
	"context"
	"payments/internal/domain/payment"

	"github.com/google/uuid"
)

type CreatePayment struct {
	OrderUUID string
	Amount    int
}

func (h Handler) CreatePayment(ctx context.Context, cmd CreatePayment) (pmnt *payment.Payment, err error) {
	pmnt, err = payment.NewPayment(uuid.NewString(), cmd.OrderUUID, cmd.Amount)
	if err != nil {
		return
	}

	_, err = h.paymentsRepository.AddPayment(ctx, pmnt)
	if err != nil {
		return
	}
	return
}
