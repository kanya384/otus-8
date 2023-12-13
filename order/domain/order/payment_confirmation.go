package order

func (o *Order) PaymentConfirmed() {
	o.paymentChecked = true
}

func (o *Order) PaymentNotConfirmed() {
	o.paymentChecked = false
	o.failed = true
}
