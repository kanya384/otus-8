package saga

import (
	"order/pkg/event"
)

type OrderInitialized_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
}

type OrderPaymentConfirmed_v1 struct {
	Header      event.EventHeader `json:"header"`
	PaymentUUID string            `json:"payment_uuid"`
	OrderUUID   string            `json:"order_uuid"`
}

type PaymentRefund_v1 struct {
	Header      event.EventHeader `json:"header"`
	PaymentUUID string            `json:"payment_uuid"`
}

type OrderPaymentFailed_v1 struct {
	Header      event.EventHeader `json:"header"`
	PaymentUUID string            `json:"payment_uuid"`
	OrderUUID   string            `json:"order_uuid"`
	Reason      string            `json:"reason"`
}

type StockProductsReserved_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
}

type StockProductsCancelReserve_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
	Items     []struct{}        `json:"items"`
}

type StockProductsReserveFailed_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
	Reason    string            `json:"reason"`
}

type DeliveryCourierReserved_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
}

type DeliveryCourierReserveFailed_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
}

type CheckOrderPayment_v1 struct {
	Header      event.EventHeader `json:"header"`
	PaymentUUID string            `json:"payment_uuid"`
	OrderUUID   string            `json:"order_uuid"`
}

type ReserveOrderProducts_v1 struct {
	Header    event.EventHeader `json:"header"`
	OrderUUID string            `json:"order_uuid"`
	Items     []struct{}        `json:"items"`
}

type ReserveDeliveryManForOrder_v1 struct {
	Header       event.EventHeader `json:"header"`
	OrderUUID    string            `json:"order_uuid"`
	CustomerName string            `json:"customer_name"`
	Address      string            `json:"address"`
	StartTime    string            `json:"start_time"`
	EndTime      string            `json:"end_time"`
}
