package server

import (
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
	userSchema := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);
	`

	productSchema := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);
	`

	seedProducts := `
	INSERT INTO products (name, price) VALUES 
		('Laptop', 999.99),
		('Smartphone', 499.99),
		('Headphones', 199.99)
	ON CONFLICT DO NOTHING;
	`

	if _, err := db.Exec(userSchema); err != nil {
		return err
	}

	if _, err := db.Exec(productSchema); err != nil {
		return err
	}

	if _, err := db.Exec(seedProducts); err != nil {
		return err
	}

	return nil
}
