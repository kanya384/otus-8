package command

import (
	"context"
	"order/internal/domain/order"
	"order/pkg/dateTime"

	"github.com/google/uuid"
)

type CreateOrder struct {
	CustomerName                string
	DeliveryAddress             string
	OrderItems                  []OrderItem
	PaymentUUID                 string
	ComfortaleDeliveryTimeStart string
	ComfortaleDeliveryTimeEnd   string
}

type OrderItem struct {
	ProductUUID string
	Price       int
	Quantity    int
}

func (h Handler) CreateOrder(ctx context.Context, cmd CreateOrder) (or *order.Order, err error) {
	startTime, err := dateTime.NewDateTimeFromString(cmd.ComfortaleDeliveryTimeStart)
	if err != nil {
		return
	}

	endTime, err := dateTime.NewDateTimeFromString(cmd.ComfortaleDeliveryTimeEnd)
	if err != nil {
		return
	}

	orderItems := make([]order.OrderItem, 0, len(cmd.OrderItems))

	for _, oi := range cmd.OrderItems {
		orderItems = append(orderItems, order.OrderItem{
			ProductUUID: oi.ProductUUID,
			Quantity:    oi.Quantity,
			Price:       oi.Price,
		})
	}

	or, err = order.NewOrder(uuid.NewString(), cmd.CustomerName, cmd.DeliveryAddress, orderItems, cmd.PaymentUUID, startTime, endTime)
	if err != nil {
		return
	}

	_, err = h.ordersRepository.AddOrder(ctx, or)
	if err != nil {
		return
	}
	return
}
