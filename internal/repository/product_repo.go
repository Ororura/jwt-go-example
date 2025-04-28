package repository

import (
	"jwt-go/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ProductSQLiteRepository struct {
	db *sqlx.DB
}

func NewProductSQLiteRepository(db *sqlx.DB) *ProductSQLiteRepository {
	return &ProductSQLiteRepository{db: db}
}

func (r *ProductSQLiteRepository) ListProducts() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Select(&products, "SELECT id, name, price FROM products")
	return products, err
}
