package jwt_provider

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type JwtProvider struct {
	secret string
}

func New(secret string) JwtProvider {
	return JwtProvider{secret: secret}
}

func (jwtProvider JwtProvider) Sign(id string) (string, error) {
	claims := Payload{
		id,
		jwt.RegisteredClaims{
			Issuer:    "post_board",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtProvider.secret)
}

func (jwtProvider JwtProvider) Parse(tokenString string) (Payload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtProvider.secret, nil
	})

	if err != nil {
		return Payload{}, err
	}

	claims, ok := token.Claims.(Payload)
	if !ok {
		return Payload{}, errors.New("unknown jwt payload")
	}
	return claims, nil
}
