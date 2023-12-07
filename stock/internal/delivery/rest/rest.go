package rest

import (
	"context"
	"fmt"
	"net/http"
	"stock/internal/domain"

	"github.com/google/uuid"
)

type Rest struct {
	port     int
	srv      *http.Server
	services Service
}

func New(port int, services Service) *Rest {

	var d = &Rest{
		port: port,
	}

	d.srv = d.initRouter(port)
	d.services = services
	return d
}

func (d *Rest) Run() (err error) {
	fmt.Printf("\nServer started at: http://localhost:%d\n", d.port)
	err = d.srv.ListenAndServe()
	return
}

func (d *Rest) Stop(ctx context.Context) (err error) {
	err = d.srv.Shutdown(ctx)
	return
}

type Service interface {
	CreateProduct(ctx context.Context, product domain.Product) (err error)
	ReadProducts(ctx context.Context) (products []domain.Product, err error)
	ReadProductById(ctx context.Context, productId uuid.UUID) (product domain.Product, err error)
}
