package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ProductId  uuid.UUID `json:"product_id" db:"product_id"`
	Name       string    `json:"name" db:"name"`
	Quantity   int       `json:"quantity" db:"quantity"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}

type ProductReserveItem struct {
	ProductId uuid.UUID
	Quantity  int
}
