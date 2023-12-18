package repository

import (
	"context"
	"database/sql"
	"fmt"
	"payments/internal/domain/payment"

	"github.com/jmoiron/sqlx"
)

type PaymentModel struct {
	UUID string `db:"uuid"`

	OrderUUID string `db:"order_uuid"`

	Amount  int  `db:"amount"`
	Success bool `db:"success"`

	Failed bool `db:"failed"`
}

type PaymentPostgresRepository struct {
	db *sqlx.DB
}

func NewPaymentPostgresRepository(db *sqlx.DB) PaymentPostgresRepository {
	return PaymentPostgresRepository{
		db: db,
	}
}

func (r PaymentPostgresRepository) AddPayment(ctx context.Context, or *payment.Payment) (*payment.Payment, error) {
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

		paymentModel, err := r.marshalPayment(payment)
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
	paymentModel := PaymentModel{}
	err := db.GetContext(ctx, &paymentModel, `
		SELECT 
			* 
		FROM 
			payments
		WHERE
			uuid = $1
	`, paymentUUID)

	if err != nil {
		return &payment.Payment{}, fmt.Errorf("could not get payment: %w", err)
	}
	or, err := payment.UnmarshalPaymentFromDatabase(paymentModel.UUID, paymentModel.OrderUUID, paymentModel.Amount, paymentModel.Success, paymentModel.Failed)
	if err != nil {
		return &payment.Payment{}, fmt.Errorf("could not unmarshal payment: %w", err)
	}
	return or, nil
}

func (r PaymentPostgresRepository) marshalPayment(pmnt *payment.Payment) (PaymentModel, error) {
	paymentModel := PaymentModel{
		UUID: pmnt.Uuid(),

		OrderUUID: pmnt.OrderUUID(),
		Amount:    pmnt.Amount(),

		Success: pmnt.Success(),
		Failed:  pmnt.Failed(),
	}

	return paymentModel, nil
}
