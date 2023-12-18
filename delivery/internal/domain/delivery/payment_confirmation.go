package payment

func (o *Delivery) DeliveryProcessed() {
	o.processed = true
}

func (o *Delivery) DeliveryFail() {
	o.fail = true
	o.processed = false
}
