package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func InitializeDatabaseSchema(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS product (
			product_id UUID PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			modified_at TIMESTAMP NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("could not initialize database schema: %w", err)
	}

	return nil
}
