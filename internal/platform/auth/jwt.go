package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID     uint   `json:"user_id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	EmployeeID *uint  `json:"employee_id,omitempty"`
	jwt.RegisteredClaims
}

type TokenProvider struct {
	secret string
	ttl    time.Duration
}

func NewTokenProvider(secret string, ttl time.Duration) *TokenProvider {
	return &TokenProvider{
		secret: secret,
		ttl:    ttl,
	}
}

func (p *TokenProvider) GenerateToken(userID uint, email, role string, employeeID *uint) (string, error) {
	claims := &Claims{
		UserID:     userID,
		Email:      email,
		Role:       role,
		EmployeeID: employeeID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secret))
}

func (p *TokenProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(p.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
