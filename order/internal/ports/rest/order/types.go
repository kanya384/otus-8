package order

type CreateOrderRequest struct {
	CustomerName                string      `json:"customer_name" binding:"required" example:"Vasiliy Petrov"`
	DeliveryAddress             string      `json:"delivery_address" binding:"required" example:"st. Krasnaya 10"`
	OrderItems                  []OrderItem `json:"order_items" binding:"required"`
	PaymentUUID                 string      `json:"payment_uuid" binding:"required" example:"f9d62750-9d9d-11ee-8c90-0242ac120002"`
	ComfortaleDeliveryTimeStart string      `json:"comfortale_delivery_time_start" binding:"required" example:"10:00"`
	ComfortaleDeliveryTimeEnd   string      `json:"comfortale_delivery_time_end" binding:"required" example:"11:30"`
}

type OrderItem struct {
	ProductUUID string `json:"product_uuid" example:"0194e35a-9d9e-11ee-8c90-0242ac120002"`
	Price       int    `json:"price" example:"123"`
	Quantity    int    `json:"quantity" example:"2"`
}

type OrderResponse struct {
	UUID string `json:"uuid"`

	CustomerName         string      `json:"customer_name"`
	OrderItems           []OrderItem `json:"order_items"`
	StockReservationDone bool        `json:"stock_reservation_done"`
	PaymentUUID          string      `json:"payment_uuid"`
	PaymentChecked       bool        `json:"payment_checked"`

	DeliveryAddress             string `json:"delivery_address"`
	ComfortaleDeliveryTimeStart string `json:"comfortale_delivery_time_start" binding:"required"`
	ComfortaleDeliveryTimeEnd   string `json:"comfortale_delivery_time_end" binding:"required"`
	DeliverySlotReserved        bool   `json:"delivery_slot_reserved"`

	Finalized bool `json:"finalized"`
	Failed    bool `json:"failed"`
}
