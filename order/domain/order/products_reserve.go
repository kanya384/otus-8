package order

func (o *Order) ProductsReservedOnStock() {
	o.stockReservationDone = true
}

func (o *Order) ProductsReservedFailedOnStock() {
	o.stockReservationDone = false
	o.failed = true
}
