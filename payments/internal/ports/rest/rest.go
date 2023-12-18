package rest

import (
	"context"
	"fmt"
	"net/http"
	"payments/internal/app/command"
	"payments/internal/domain/payment"
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
	fmt.Printf("started http server at: %d\n", d.port)
	return d.srv.ListenAndServe()
}

func (d *Rest) Stop(ctx context.Context) (err error) {
	return d.srv.Shutdown(ctx)
}

type CommandHandler interface {
	CreatePayment(ctx context.Context, cmd command.CreatePayment) (order *payment.Payment, err error)
	ReadPaymentById(ctx context.Context, cmd command.ReadPayment) (order *payment.Payment, err error)
}
