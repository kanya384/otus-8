package payment

func (o *Payment) PaymentSuccess() {
	o.failed = false
	o.success = true
}

func (o *Payment) PaymentFailed() {
	o.failed = true
	o.success = false
}
