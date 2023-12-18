package payment

type CreatePaymentRequest struct {
	OrderUUID string `json:"order_uuid" binding:"required" example:"f9d62750-9d9d-11ee-8c90-0242ac120002"`
	Amount    int    `json:"amount" binding:"required" example:"122"`
}

type PaymentResponse struct {
	UUID string `json:"uuid"`

	OrderUUID string `json:"order_uuid"`

	Amount int `json:"amount" example:"122"`

	Success bool `json:"success" example:"true"`
	Failed  bool `json:"failde" example:"false"`
}
