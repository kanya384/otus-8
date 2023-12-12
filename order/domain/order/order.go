package order

import (
	"encoding/json"
	"errors"
	"order/pkg/dateTime"
)

type Order struct {
	uuid string

	customerName string

	orderItems           []OrderItem
	stockReservationDone bool

	paymentUUID    string
	paymentChecked bool

	deliveryAddress             string
	comfortaleDeliveryTimeStart dateTime.DateTime
	comfortaleDeliveryTimeEnd   dateTime.DateTime
	deliverySlotReserved        bool

	finalized bool
	failed    bool
}

func NewOrder(uuid string, customerName, deliveryAddress string, orderItems []OrderItem, paymentUUID string, comfortaleDeliveryTimeStart dateTime.DateTime, comfortaleDeliveryTimeEnd dateTime.DateTime) (*Order, error) {
	if uuid == "" {
		return nil, errors.New("empty order uuid")
	}

	if customerName == "" {
		return nil, errors.New("empty customer name")
	}

	if deliveryAddress == "" {
		return nil, errors.New("empty delivery address")
	}

	if deliveryAddress == "" {
		return nil, errors.New("empty delivery address")
	}

	if len(orderItems) == 0 {
		return nil, errors.New("empty ordered items")
	}

	if deliveryAddress == "" {
		return nil, errors.New("empty delivery address")
	}

	if paymentUUID == "" {
		return nil, errors.New("paymentUUID empty")
	}

	if comfortaleDeliveryTimeStart.IsNil() {
		return nil, errors.New("nil comfortable delivery time start")
	}

	if comfortaleDeliveryTimeEnd.IsNil() {
		return nil, errors.New("nil comfortable delivery time end")
	}

	if comfortaleDeliveryTimeStart.After(comfortaleDeliveryTimeEnd) {
		return nil, errors.New("time start is after time end")
	}

	return &Order{
		uuid:                        uuid,
		customerName:                customerName,
		orderItems:                  orderItems,
		stockReservationDone:        false,
		paymentUUID:                 paymentUUID,
		paymentChecked:              false,
		deliveryAddress:             deliveryAddress,
		comfortaleDeliveryTimeStart: comfortaleDeliveryTimeStart,
		comfortaleDeliveryTimeEnd:   comfortaleDeliveryTimeEnd,
		deliverySlotReserved:        false,
		finalized:                   false,
		failed:                      false,
	}, nil
}

func (o Order) Uuid() string {
	return o.uuid
}

func (o Order) CustomerName() string {
	return o.customerName
}

func (o Order) OrderItems() []OrderItem {
	return o.orderItems
}

func (o Order) StockReservationDone() bool {
	return o.stockReservationDone
}

func (o Order) PaymentUUID() string {
	return o.paymentUUID
}

func (o Order) PaymentChecked() bool {
	return o.paymentChecked
}

func (o Order) DeliveryAddress() string {
	return o.deliveryAddress
}

func (o Order) ComfortaleDeliveryTimeStart() dateTime.DateTime {
	return o.comfortaleDeliveryTimeStart
}

func (o Order) ComfortaleDeliveryTimeEnd() dateTime.DateTime {
	return o.comfortaleDeliveryTimeEnd
}

func (o Order) DeliverySlotReserved() bool {
	return o.deliverySlotReserved
}

func (o Order) Finalized() bool {
	return o.finalized
}

func (o Order) Failed() bool {
	return o.failed
}

func UnmarshalOrderFromDatabase(
	UUID string,

	CustomerName string,

	OrderItems string,
	StockReservationDone bool,

	PaymentUUID string,
	PaymentChecked bool,

	DeliveryAddress string,
	ComfortaleDeliveryTimeStart string,
	ComfortaleDeliveryTimeEnd string,
	DeliverySlotReserved bool,

	Finalized bool,
	Failed bool,
) (*Order, error) {
	var orderItems []OrderItem
	err := json.Unmarshal([]byte(OrderItems), &orderItems)
	if err != nil {
		return nil, err
	}
	startTime, err := dateTime.NewDateTimeFromString(ComfortaleDeliveryTimeStart)
	if err != nil {
		return nil, err
	}

	endTime, err := dateTime.NewDateTimeFromString(ComfortaleDeliveryTimeEnd)
	if err != nil {
		return nil, err
	}
	order, err := NewOrder(UUID, CustomerName, DeliveryAddress, orderItems, PaymentUUID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	return order, nil
}
