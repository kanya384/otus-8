package payment

import "payments/internal/domain/payment"

func ToPaymentResponse(pmnt *payment.Payment) PaymentResponse {
	return PaymentResponse{
		UUID:      pmnt.Uuid(),
		OrderUUID: pmnt.OrderUUID(),
		Amount:    pmnt.Amount(),
		Success:   pmnt.Success(),
		Failed:    pmnt.Failed(),
	}
}
