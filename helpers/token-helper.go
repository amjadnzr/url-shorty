package helpers

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenHelper struct {
	secret    string
	publisher string
}

func NewTokenHelper(jwtSecret string) *TokenHelper {
	return &TokenHelper{
		secret:    jwtSecret,
		publisher: "url-shortly",
	}
}

type Claims struct {
	jwt.RegisteredClaims
}

func (t *TokenHelper) GenerateJWTToken(userId int64) (string, error) {
	claims := &Claims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    t.publisher,
			Subject:   fmt.Sprintf("user-%d", userId),

			Audience: []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte(t.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (t *TokenHelper) ValidateJWTToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
			return []byte(t.secret), nil
		}
		return nil, errors.New("token signing invalid")
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("Claim could not be decoded")
	}

	return claims, nil
}
