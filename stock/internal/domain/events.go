package domain

import (
	"stock/pkg/types"

	"github.com/google/uuid"
)

type ReserveProducts_v1 struct {
	Header   types.EventHeader    `json:"header"`
	OrderId  uuid.UUID            `json:"order_id"`
	Products []ReserveProductItem `json:"products"`
}

type ReserveProductItem struct {
	ProductId uuid.UUID `json:"product_id"`
	Quantity  int       `json:"quantity"`
}

type ReserveProductsSucces_v1 struct {
	Header  types.EventHeader `json:"header"`
	OrderId uuid.UUID         `json:"order_id"`
}

type ReserveProductsFail_v1 struct {
	Header  types.EventHeader `json:"header"`
	OrderId uuid.UUID         `json:"order_id"`
	Reason  string            `json:"reason"`
}

type CancelReserveProducts_v1 struct {
	Header   types.EventHeader    `json:"header"`
	OrderId  uuid.UUID            `json:"order_id"`
	Products []ReserveProductItem `json:"products"`
}

type CancelReserveProductsSuccess_v1 struct {
	Header  types.EventHeader `json:"header"`
	OrderId uuid.UUID         `json:"order_id"`
}
