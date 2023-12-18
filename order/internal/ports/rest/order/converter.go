package order

import (
	"order/internal/domain/order"
)

func ToOrderResponse(do *order.Order) OrderResponse {
	orderItems := make([]OrderItem, 0, len(do.OrderItems()))

	for _, oi := range do.OrderItems() {
		orderItems = append(orderItems, OrderItem{
			ProductUUID: oi.ProductUUID,
			Quantity:    oi.Quantity,
			Price:       oi.Price,
		})
	}
	return OrderResponse{
		UUID:                        do.Uuid(),
		CustomerName:                do.CustomerName(),
		OrderItems:                  orderItems,
		StockReservationDone:        do.StockReservationDone(),
		PaymentUUID:                 do.PaymentUUID(),
		PaymentChecked:              do.PaymentChecked(),
		DeliveryAddress:             do.DeliveryAddress(),
		ComfortaleDeliveryTimeStart: do.ComfortaleDeliveryTimeEnd().ToString(),
		ComfortaleDeliveryTimeEnd:   do.ComfortaleDeliveryTimeEnd().ToString(),
		DeliverySlotReserved:        do.DeliverySlotReserved(),
		Finalized:                   do.Finalized(),
		Failed:                      do.Failed(),
	}
}
