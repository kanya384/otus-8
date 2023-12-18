package rest

import (
	"context"
	"net/http"

	"payments/internal/app/command"
	paymentsRequest "payments/internal/ports/rest/payment"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreatePayment
// @Summary создать платеж.
// @Description создать платеж.
// @Tags orders
// @Accept  json
// @Produce json
// @Param   data 		body 		orderRequest.CreateOrderRequest 		true  "Данные для создания продукта"
// @Success 200			{object}    orderRequest.OrderResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/order [post]
func (d *Rest) CreatePayment(c *gin.Context) {
	request := paymentsRequest.CreatePaymentRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	order, err := d.commandHandlers.CreatePayment(context.Background(), command.CreatePayment{
		OrderUUID: request.OrderUUID,
		Amount:    request.Amount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, paymentsRequest.ToPaymentResponse(order))
}

// ReadOrderById
// @Summary получить платеж по id.
// @Description получить платеж по id.
// @Tags orders
// @Accept  json
// @Produce json
// @Param   id 			path 		string 						true  "Идентификатор платежа"
// @Success 200			{object}    orderRequest.OrderResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/order/{id} [get]
func (d *Rest) ReadPaymentById(c *gin.Context) {

	orderUUID := c.Param("id")

	_, err := uuid.Parse(orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	product, err := d.commandHandlers.ReadPaymentById(context.Background(), command.ReadPayment{
		PaymentUUID: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
