package rest

import (
	"context"
	"net/http"
	"time"

	productRequest "stock/internal/delivery/rest/product"

	"stock/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProduct
// @Summary создать продукт.
// @Description создать продукт.
// @Tags products
// @Accept  json
// @Produce json
// @Param   data 		body 		productRequest.CreateProductRequest 		true  "Данные для создания продукта"
// @Success 200			{object}    domain.Product
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/product [post]
func (d *Rest) CreateProduct(c *gin.Context) {
	request := productRequest.CreateProductRequest{}

	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := domain.Product{
		ProductId:  uuid.New(),
		Name:       request.Name,
		Quantity:   request.Quantity,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err := d.services.CreateProduct(context.Background(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// ReadProducts
// @Summary получить список продуктов.
// @Description получить список продуктов.
// @Tags products
// @Accept  json
// @Produce json
// @Success 200			{array}    	domain.Product
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/product [get]
func (d *Rest) ReadProducts(c *gin.Context) {
	products, err := d.services.ReadProducts(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// ReadProductById
// @Summary получить продукт по id.
// @Description получить продукт по id.
// @Tags products
// @Accept  json
// @Produce json
// @Param   id 			path 		string 						true  "Идентификатор продукта"
// @Success 200			{object}    domain.Product
// @Failure 400 		{object}    ErrorResponse
// @Failure 401 		{object}    ErrorResponse
// @Router /api/product/{id} [get]
func (d *Rest) ReadProductById(c *gin.Context) {
	productId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := d.services.ReadProductById(context.Background(), productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}
