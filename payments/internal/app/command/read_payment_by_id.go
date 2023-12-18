package command

import (
	"context"
	"payments/internal/domain/payment"
)

type ReadPayment struct {
	PaymentUUID string
}

func (h Handler) ReadPaymentById(ctx context.Context, cmd ReadPayment) (pmnt *payment.Payment, err error) {

	pmnt, err = h.paymentsRepository.ReadPayment(ctx, cmd.PaymentUUID)
	if err != nil {
		return
	}
	return
}
