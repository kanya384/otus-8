package order

func (o *Order) PaymentComfirmed() {
	o.paymentChecked = true
}

func (o *Order) PaymentNotComfirmed() {
	o.paymentChecked = false
	o.failed = true
}
