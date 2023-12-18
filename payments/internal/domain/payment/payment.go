package payment

import (
	"errors"
)

type Payment struct {
	uuid string

	orderUUID string

	amount int

	success bool
	failed  bool
}

func NewPayment(
	uuid string,
	orderUUID string,
	amount int,
) (*Payment, error) {
	if uuid == "" {
		return nil, errors.New("empty payment uuid")
	}

	if orderUUID == "" {
		return nil, errors.New("empty order uuid")
	}

	if amount <= 0 {
		return nil, errors.New("amount less than zero")
	}

	return &Payment{
		uuid:      uuid,
		orderUUID: orderUUID,
		amount:    amount,
	}, nil
}

func (o Payment) Uuid() string {
	return o.uuid
}

func (o Payment) OrderUUID() string {
	return o.orderUUID
}

func (o Payment) Amount() int {
	return o.amount
}

func (o Payment) Success() bool {
	return o.success
}

func (o Payment) Failed() bool {
	return o.failed
}

func UnmarshalPaymentFromDatabase(
	UUID string,

	OrderUUID string,

	Amount int,

	success bool,

	failed bool,

) (*Payment, error) {

	return &Payment{
		uuid:      UUID,
		orderUUID: OrderUUID,
		amount:    Amount,
		success:   success,
		failed:    failed,
	}, nil
}
