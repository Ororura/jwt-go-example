package domain

type Product struct {
	ID    int     `db:"id" json:"id"`
	Name  string  `db:"name" json:"name"`
	Price float64 `db:"price" json:"price"`
}

type ProductRepository interface {
	ListProducts() ([]Product, error)
}
