package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitializeDatabaseSchema(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS courier (
			uuid UUID PRIMARY KEY,
			name varchar NOT NULL
		);

		CREATE TABLE IF NOT EXISTS delivery (
			uuid UUID PRIMARY KEY,
			courier_uuid UUID NOT NULL,
			address varchar not null,
			processed boolean not null
		);
	`)
	if err != nil {
		return fmt.Errorf("could not initialize database schema: %w", err)
	}

	return nil
}
