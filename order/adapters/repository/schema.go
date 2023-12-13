package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitializeDatabaseSchema(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			uuid UUID PRIMARY KEY,
			customer_name VARCHAR(255) NOT NULL,
			order_items VARCHAR NOT NULL,
			stock_reservation_done boolean not null,
			payment_uuid UUID not null,
			payment_checked boolean not null,
			delivery_address VARCHAR not null,
			comfortable_delivery_time_start VARCHAR(6) NOT NULL,
			comfortable_delivery_time_end VARCHAR(6) NOT NULL,
			delivery_slot_reserved boolean not null,
			finalized boolean not null,
			failed boolean not null
		);
	`)
	if err != nil {
		return fmt.Errorf("could not initialize database schema: %w", err)
	}

	return nil
}
