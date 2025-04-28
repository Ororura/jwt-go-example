package usecase

import (
	"errors"
	"jwt-go/internal/domain"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo      domain.UserRepository
	jwtSecret []byte
}

func NewAuthUsecase(r domain.UserRepository, secret string) *AuthUsecase {
	return &AuthUsecase{
		repo:      r,
		jwtSecret: []byte(secret),
	}
}

func (u *AuthUsecase) Register(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password: " + err.Error())
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return u.repo.Save(user)
}

func (u *AuthUsecase) Login(username, password string) (string, error) {
	user, err := u.repo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	claims := jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}

func (u *AuthUsecase) ValidateToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return u.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
