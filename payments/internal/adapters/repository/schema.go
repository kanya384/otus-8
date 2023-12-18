package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitializeDatabaseSchema(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			uuid UUID PRIMARY KEY,
			order_uuid UUID NOT NULL,
			amount integer NOT NULL,
			success boolean not null,
			failed boolean not null
		);
	`)
	if err != nil {
		return fmt.Errorf("could not initialize database schema: %w", err)
	}

	return nil
}
