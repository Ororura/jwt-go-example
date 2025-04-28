package domain

type User struct {
	Username string
	Password string
}

type UserRepository interface {
	Save(user *User) error
	GetByUsername(username string) (*User, error)
}
