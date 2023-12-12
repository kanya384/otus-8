package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"order/domain/order"

	"github.com/jmoiron/sqlx"
)

type OrderModel struct {
	UUID string `db:"uuid"`

	CustomerName string `db:"customer_name"`

	OrderItems           string `db:"order_items"`
	StockReservationDone bool   `db:"stock_reservation_done"`

	PaymentUUID    string `db:"payment_uuid"`
	PaymentChecked bool   `db:"payment_checked"`

	DeliveryAddress             string `db:"payment_uuid"`
	ComfortaleDeliveryTimeStart string `db:"comfortable_delivery_time_start"`
	ComfortaleDeliveryTimeEnd   string `db:"comfortable_delivery_time_end"`
	DeliverySlotReserved        bool   `db:"delivery_slot_reserved"`

	Finalized bool `db:"finalized"`
	Failed    bool `db:"failed"`
}

type OrderPostgresRepository struct {
	db *sqlx.DB
}

func NewOrderPostgresRepository(db *sqlx.DB) OrderPostgresRepository {
	return OrderPostgresRepository{
		db: db,
	}
}

func (r OrderPostgresRepository) AddOrder(ctx context.Context, or *order.Order) error {
	orderModel, err := r.marshalOrder(or)
	if err != nil {
		return err
	}
	_, err = r.db.NamedExecContext(
		ctx,
		`
		INSERT INTO 
    		orders (uuid, customer_name, order_items, stock_reservation_done, payment_uuid, payment_checked, payment_uuid, comfortable_delivery_time_start, comfortable_delivery_time_end, delivery_slot_reserved, finalized, failed) 
		VALUES 
		    (:uuid, :customer_name, :order_items, :stock_reservation_done, :payment_uuid, :payment_checked, :payment_uuid, :comfortable_delivery_time_start, :comfortable_delivery_time_end, :delivery_slot_reserved, :finalized, :failed) 
		ON CONFLICT DO NOTHING`,
		orderModel,
	)
	if err != nil {
		return fmt.Errorf("could not save order: %w", err)
	}
	return nil
}

func (r OrderPostgresRepository) ReadOrder(ctx context.Context, orderUUID string) (*order.Order, error) {
	orderModel := OrderModel{}
	err := r.db.GetContext(ctx, &orderModel, `
		SELECT 
			* 
		FROM 
			orders
		WHERE
			uuid = $1
	`, orderUUID)

	if err != nil {
		return &order.Order{}, fmt.Errorf("could not get order: %w", err)
	}
	or, err := order.UnmarshalOrderFromDatabase(orderModel.UUID, orderModel.CustomerName, orderModel.OrderItems, orderModel.StockReservationDone, orderModel.PaymentUUID, orderModel.PaymentChecked, orderModel.DeliveryAddress, orderModel.ComfortaleDeliveryTimeStart, orderModel.ComfortaleDeliveryTimeEnd, orderModel.DeliverySlotReserved, orderModel.Finalized, orderModel.Failed)
	if err != nil {
		return &order.Order{}, fmt.Errorf("could not unmarshal order: %w", err)
	}
	return or, nil
}

func (r OrderPostgresRepository) UpdateOrder(
	ctx context.Context,
	orderUUID string,
	updateFn func(ctx context.Context, or *order.Order) (*order.Order, error),
) error {
	return nil
}

func (r OrderPostgresRepository) marshalOrder(or *order.Order) (OrderModel, error) {
	orderModel := OrderModel{
		UUID: or.Uuid(),

		CustomerName: or.CustomerName(),

		StockReservationDone: or.StockReservationDone(),
		PaymentUUID:          or.PaymentUUID(),
		PaymentChecked:       or.PaymentChecked(),

		DeliveryAddress:             or.DeliveryAddress(),
		ComfortaleDeliveryTimeStart: or.ComfortaleDeliveryTimeStart().ToString(),
		ComfortaleDeliveryTimeEnd:   or.ComfortaleDeliveryTimeEnd().ToString(),
		DeliverySlotReserved:        or.DeliverySlotReserved(),

		Finalized: or.Finalized(),
		Failed:    or.Failed(),
	}

	products, err := json.Marshal(or.OrderItems())
	if err != nil {
		return orderModel, err
	}
	orderModel.OrderItems = string(products)
	return orderModel, nil
}
