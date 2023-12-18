package rest

import (
	"context"
	"fmt"
	"net/http"
	"order/internal/app/command"
	"order/internal/domain/order"
)

type Rest struct {
	port            int
	srv             *http.Server
	commandHandlers CommandHandler
}

func New(port int, commandHandlers CommandHandler) *Rest {

	var d = &Rest{
		port: port,
	}

	d.srv = d.initRouter(port)
	d.commandHandlers = commandHandlers
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

type CommandHandler interface {
	CreateOrder(ctx context.Context, cmd command.CreateOrder) (order *order.Order, err error)
	ReadOrderById(ctx context.Context, cmd command.ReadOrder) (order *order.Order, err error)
}
