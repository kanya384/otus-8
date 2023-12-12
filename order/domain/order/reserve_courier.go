package order

func (o *Order) CourierReserved() {
	o.deliverySlotReserved = true
	o.finalized = true
}

func (o *Order) CourierReserveFailed() {
	o.deliverySlotReserved = false
	o.failed = true
}
