package repository

import (
	"context"
	"database/sql"
	"fmt"
	"stock/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	if db == nil {
		panic("db is nil")
	}

	return ProductRepository{db: db}
}

func (s ProductRepository) AddProduct(ctx context.Context, product domain.Product) error {
	_, err := s.db.NamedExecContext(ctx, `
		INSERT INTO 
		    Product (product_id, name, quantity, created_at, modified_at) 
		VALUES (:product_id, :name, :quantity, :created_at, :modified_at)
		`, product)
	if err != nil {
		return fmt.Errorf("could not add Product: %w", err)
	}

	return nil
}

func (s ProductRepository) AllProduct(ctx context.Context) ([]domain.Product, error) {
	var Product []domain.Product
	err := s.db.SelectContext(ctx, &Product, `
		SELECT 
		    * 
		FROM 
		    product
	`)
	if err != nil {
		return nil, fmt.Errorf("could not get Product: %w", err)
	}

	return Product, nil
}

func (s ProductRepository) ProductByID(ctx context.Context, productID uuid.UUID) (domain.Product, error) {
	var Product domain.Product
	err := s.db.GetContext(ctx, &Product, `
		SELECT 
		    * 
		FROM 
		    product
		WHERE
		    product_id = $1
	`, productID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("could not get Product: %w", err)
	}

	return Product, nil
}

func (s ProductRepository) UpdateByID(ctx context.Context, productID uuid.UUID, updateFn func(product domain.Product) (domain.Product, error)) (domain.Product, error) {
	var pr domain.Product

	err := updateInTx(ctx, s.db, sql.LevelSerializable, func(ctx context.Context, tx *sqlx.Tx) error {
		var err error
		pr, err = s.oneProductByIDInTx(ctx, productID, tx)
		if err != nil {
			return err
		}

		pr, err = updateFn(pr)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE product SET name = $1, quantity = $2, modified_at = NOW() WHERE vip_bundle_id = $3
		`, pr.Name, pr.Quantity, pr.ProductId)

		if err != nil {
			return fmt.Errorf("could not update product: %w", err)
		}

		return nil
	})
	if err != nil {
		return domain.Product{}, fmt.Errorf("could not update product: %w", err)
	}

	return pr, nil
}

func (s ProductRepository) oneProductByIDInTx(ctx context.Context, productID uuid.UUID, tx *sqlx.Tx) (domain.Product, error) {
	var Product domain.Product
	err := tx.GetContext(ctx, &Product, `
		SELECT 
		    * 
		FROM 
		    product
		WHERE
		    product_id = $1
	`, productID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("could not get Product: %w", err)
	}

	return Product, nil
}
