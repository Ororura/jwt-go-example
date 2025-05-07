package usecase

import (
	"errors"
	"jwt-go/internal/domain"
	"time"

	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo         domain.UserRepository
	jwtSecret    []byte
	refreshStore map[string]string
	mu           sync.Mutex
}

func NewAuthUsecase(r domain.UserRepository, secret string) *AuthUsecase {
	return &AuthUsecase{
		repo:         r,
		jwtSecret:    []byte(secret),
		refreshStore: make(map[string]string),
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

func (u *AuthUsecase) Login(username, password string) (accessToken string, refreshToken string, err error) {
	user, err := u.repo.GetByUsername(username)
	if err != nil {
		return "", "", errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid username or password")
	}

	accessToken, err = u.generateAccessToken(username)
	if err != nil {
		return "", "", err
	}

	refreshToken = uuid.New().String()

	u.mu.Lock()
	u.refreshStore[refreshToken] = username
	u.mu.Unlock()

	return accessToken, refreshToken, nil
}

func (u *AuthUsecase) RefreshAccessToken(refreshToken string) (string, error) {
	u.mu.Lock()
	username, ok := u.refreshStore[refreshToken]
	u.mu.Unlock()

	if !ok {
		return "", errors.New("invalid refresh token")
	}

	return u.generateAccessToken(username)
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

func (u *AuthUsecase) generateAccessToken(username string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   username,
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "auth-app",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}
