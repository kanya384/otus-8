package rest

import (
	"context"
	"net/http"

	"order/internal/app/command"
	orderRequest "order/internal/ports/rest/order"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateOrder
// @Summary создать заказ.
// @Description создать заказ.
// @Tags orders
// @Accept  json
// @Produce json
// @Param   data 		body 		orderRequest.CreateOrderRequest 		true  "Данные для создания продукта"
// @Success 200			{object}    orderRequest.OrderResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/order [post]
func (d *Rest) CreateOrder(c *gin.Context) {
	request := orderRequest.CreateOrderRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	orderItems := make([]command.OrderItem, 0, len(request.OrderItems))

	for _, oi := range request.OrderItems {
		orderItems = append(orderItems, command.OrderItem{
			ProductUUID: oi.ProductUUID,
			Quantity:    oi.Quantity,
			Price:       oi.Price,
		})
	}

	order, err := d.commandHandlers.CreateOrder(context.Background(), command.CreateOrder{
		CustomerName:                request.CustomerName,
		DeliveryAddress:             request.DeliveryAddress,
		OrderItems:                  orderItems,
		PaymentUUID:                 request.PaymentUUID,
		ComfortaleDeliveryTimeStart: request.ComfortaleDeliveryTimeStart,
		ComfortaleDeliveryTimeEnd:   request.ComfortaleDeliveryTimeEnd,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, orderRequest.ToOrderResponse(order))
}

// ReadOrderById
// @Summary получить заказ по id.
// @Description получить заказ по id.
// @Tags orders
// @Accept  json
// @Produce json
// @Param   id 			path 		string 						true  "Идентификатор заказа"
// @Success 200			{object}    orderRequest.OrderResponse
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/order/{id} [get]
func (d *Rest) ReadOrderById(c *gin.Context) {

	orderUUID := c.Param("id")

	_, err := uuid.Parse(orderUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	product, err := d.commandHandlers.ReadOrderById(context.Background(), command.ReadOrder{
		OrderUUID: c.Param("id"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
