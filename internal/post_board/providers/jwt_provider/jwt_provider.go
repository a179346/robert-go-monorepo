package jwt_provider

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type JwtProvider struct {
	secret        []byte
	expireSeconds int
}

func New(secret string, expireSeconds int) JwtProvider {
	return JwtProvider{
		secret:        []byte(secret),
		expireSeconds: expireSeconds,
	}
}

func (jwtProvider JwtProvider) Sign(id string) (string, error) {
	claims := Claims{
		id,
		jwt.RegisteredClaims{
			Issuer:    "post-board",
			Subject:   "auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtProvider.expireSeconds) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtProvider.secret)
}

func (jwtProvider JwtProvider) Parse(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtProvider.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("unknown jwt payload")
	}
	return claims, nil
}
