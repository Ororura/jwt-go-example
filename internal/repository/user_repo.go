package repository

import (
	"errors"
	"jwt-go/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserSQLiteRepository struct {
	db *sqlx.DB
}

func NewUserSQLiteRepository(db *sqlx.DB) *UserSQLiteRepository {
	return &UserSQLiteRepository{db: db}
}

func (r *UserSQLiteRepository) Save(user *domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, user.Password)
	if err != nil {
		return errors.New("failed to save user: " + err.Error())
	}
	return nil
}

func (r *UserSQLiteRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Get(&user, "SELECT username, password FROM users WHERE username = ?", username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
