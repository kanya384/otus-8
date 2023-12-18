package payment

import (
	"errors"
)

type Delivery struct {
	uuid string

	courierUUID string

	address string

	processed bool
	fail      bool
}

func NewDelivery(
	uuid string,
	courierUUID string,
	address string,
) (*Delivery, error) {
	if uuid == "" {
		return nil, errors.New("empty payment uuid")
	}

	if courierUUID == "" {
		return nil, errors.New("empty courier uuid")
	}

	if address == "" {
		return nil, errors.New("address")
	}

	return &Delivery{
		uuid:        uuid,
		courierUUID: courierUUID,
		address:     address,
		processed:   false,
	}, nil
}

func (o Delivery) Uuid() string {
	return o.uuid
}

func (o Delivery) CourierUUID() string {
	return o.courierUUID
}

func (o Delivery) Address() string {
	return o.address
}

func (o Delivery) Processed() bool {
	return o.processed
}

func UnmarshalDeliveryFromDatabase(
	uuid string,
	courierUUID string,
	address string,
	processed bool,
) (*Delivery, error) {

	return &Delivery{
		uuid:        uuid,
		courierUUID: courierUUID,
		address:     address,
		processed:   processed,
	}, nil
}
