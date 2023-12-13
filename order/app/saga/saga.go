package saga

import (
	"context"
	"order/domain/order"
	"order/pkg/event"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type OrderSaga struct {
	commandBus *cqrs.CommandBus
	eventBus   *cqrs.EventBus
	repository order.Repository
}

func NewOrderSaga(
	commandBus *cqrs.CommandBus,
	eventBus *cqrs.EventBus,
	repository order.Repository,
) *OrderSaga {
	return &OrderSaga{
		commandBus: commandBus,
		eventBus:   eventBus,
		repository: repository,
	}
}

func (s OrderSaga) OnOrderInitialized(ctx context.Context, evnt *OrderInitialized_v1) error {
	order, err := s.repository.ReadOrder(ctx, evnt.OrderUUID)
	if err != nil {
		return err
	}

	return s.eventBus.Publish(ctx, CheckOrderPayment_v1{
		Header:      event.NewEventHeader(),
		OrderUUID:   evnt.OrderUUID,
		PaymentUUID: order.PaymentUUID(),
	})
}

func (s OrderSaga) OnPaymentConfirmed(ctx context.Context, evnt *OrderPaymentConfirmed_v1) error {
	err := s.repository.UpdateOrder(ctx, evnt.OrderUUID, func(ctx context.Context, or *order.Order) (*order.Order, error) {
		or.PaymentConfirmed()

		err := s.eventBus.Publish(ctx, ReserveOrderProducts_v1{
			Header:    event.NewEventHeader(),
			OrderUUID: evnt.OrderUUID,
		})

		return or, err
	})
	if err != nil {
		return err
	}

	return err
}

func (s OrderSaga) OnPaymentFailed(ctx context.Context, event *OrderPaymentFailed_v1) error {
	err := s.repository.UpdateOrder(ctx, event.OrderUUID, func(ctx context.Context, or *order.Order) (*order.Order, error) {
		or.PaymentNotConfirmed()
		err := s.rollbackProcess(ctx, event.OrderUUID)
		return or, err
	})
	if err != nil {
		return err
	}

	return err
}

func (s OrderSaga) OnOrderProductsReserved(ctx context.Context, evt *StockProductsReserved_v1) error {
	err := s.repository.UpdateOrder(ctx, evt.OrderUUID, func(ctx context.Context, or *order.Order) (*order.Order, error) {
		or.ProductsReservedOnStock()

		err := s.eventBus.Publish(ctx, ReserveDeliveryManForOrder_v1{
			Header:       event.NewEventHeader(),
			OrderUUID:    or.Uuid(),
			CustomerName: or.CustomerName(),
			Address:      or.DeliveryAddress(),
			StartTime:    or.ComfortaleDeliveryTimeStart().ToString(),
			EndTime:      or.ComfortaleDeliveryTimeEnd().ToString(),
		})

		return or, err
	})
	if err != nil {
		return err
	}

	return err
}

func (s OrderSaga) OnOrderProductsReserveFaile(ctx context.Context, event *StockProductsReserveFailed_v1) error {
	err := s.repository.UpdateOrder(ctx, event.OrderUUID, func(ctx context.Context, or *order.Order) (*order.Order, error) {
		or.ProductsReservedFailedOnStock()
		err := s.rollbackProcess(ctx, event.OrderUUID)
		return or, err
	})
	if err != nil {
		return err
	}

	return err
}

func (s OrderSaga) rollbackProcess(ctx context.Context, orderUUID string) error {
	or, err := s.repository.ReadOrder(ctx, orderUUID)
	if err != nil {
		return err
	}

	if or.PaymentChecked() {
		if err := s.eventBus.Publish(ctx, PaymentRefund_v1{
			Header:      event.NewEventHeader(),
			PaymentUUID: or.PaymentUUID(),
		}); err != nil {
			return err
		}
	}

	if or.StockReservationDone() {
		if err := s.eventBus.Publish(ctx, StockProductsCancelReserve_v1{
			Header:    event.NewEventHeader(),
			OrderUUID: or.Uuid(),
			Items:     []struct{}{},
		}); err != nil {
			return err
		}
	}

	//delivery cancelation not needed - last step

	return err
}
