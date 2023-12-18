package repository

import (
	"context"
	"database/sql"
	delivery "delivery/internal/domain/delivery"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DeliveryModel struct {
	UUID string `db:"uuid"`

	CourierUUID string `db:"courier_uuid"`

	Address   string `db:"address"`
	Processed bool   `db:"processed"`
}

type PaymentPostgresRepository struct {
	db *sqlx.DB
}

func NewPaymentPostgresRepository(db *sqlx.DB) PaymentPostgresRepository {
	return PaymentPostgresRepository{
		db: db,
	}
}

func (r PaymentPostgresRepository) AddDelivery(ctx context.Context, delivery *delivery.Delivery) (*payment.Payment, error) {
	paymentModel, err := r.marshalPayment(or)
	if err != nil {
		return nil, err
	}
	_, err = r.db.NamedExecContext(
		ctx,
		`
		INSERT INTO 
    		payments (uuid, payment_uuid, amount, success, failed) 
		VALUES 
		    (:uuid, :payment_uuid, :amount, :success, :failed) 
		ON CONFLICT DO NOTHING`,
		paymentModel,
	)
	if err != nil {
		return nil, fmt.Errorf("could not save payment: %w", err)
	}
	return or, nil
}

func (r PaymentPostgresRepository) ReadPayment(ctx context.Context, paymentUUID string) (*payment.Payment, error) {
	return r.getByPaymentID(ctx, paymentUUID, r.db)
}

func (r PaymentPostgresRepository) UpdatePayment(
	ctx context.Context,
	paymentUUID string,
	updateFn func(ctx context.Context, oldPayment *payment.Payment) (*payment.Payment, error),
) (*payment.Payment, error) {
	var payment *payment.Payment
	err := updateInTx(ctx, r.db, sql.LevelSerializable, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		payment, err = r.getByPaymentID(ctx, paymentUUID, tx)
		if err != nil {
			return err
		}

		payment, err = updateFn(ctx, payment)
		if err != nil {
			return err
		}

		paymentModel, err := r.marshalDelivery(payment)
		if err != nil {
			return fmt.Errorf("could not marshal payment: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE payments SET success = $1, failed = $2 WHERE uuid = $3
		`, paymentModel.Success, paymentModel.Failed, payment.Uuid())

		if err != nil {
			return fmt.Errorf("could not update payment: %w", err)
		}

		return nil
	})
	if err != nil {
		return payment, fmt.Errorf("could not update payment: %w", err)
	}
	return payment, nil
}

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

func (v PaymentPostgresRepository) getByPaymentID(ctx context.Context, paymentUUID string, db Executor) (*payment.Payment, error) {
	paymentModel := DeliveryModel{}
	err := db.GetContext(ctx, &paymentModel, `
		SELECT 
			* 
		FROM 
			payments
		WHERE
			uuid = $1
	`, paymentUUID)

	if err != nil {
		return &delivery.Delivery{}, fmt.Errorf("could not get delivery: %w", err)
	}
	or, err := delivery.UnmarshalDeliveryFromDatabase(paymentModel.UUID, paymentModel.CourierUUID, paymentModel.Address, paymentModel.Processed)
	if err != nil {
		return &delivery.Delivery{}, fmt.Errorf("could not unmarshal delivery: %w", err)
	}
	return or, nil
}

func (r PaymentPostgresRepository) marshalDelivery(dvr *delivery.Delivery) (DeliveryModel, error) {
	paymentModel := DeliveryModel{
		UUID: dvr.Uuid(),

		CourierUUID: dvr.CourierUUID(),
		Address:     dvr.Address(),
		Processed:   dvr.Processed(),
	}

	return paymentModel, nil
}
