package auth_jwt

import (
	"errors"
	"fmt"
	"time"

	post_board_config "github.com/a179346/robert-go-monorepo/services/post_board/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func Sign(id string) (string, error) {
	jwtConfig := post_board_config.GetJwtConfig()

	claims := Claims{
		id,
		jwt.RegisteredClaims{
			Issuer:    "post-board",
			Subject:   "auth",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-10 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.ExpireSeconds) * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtConfig.Secret)
}

func Parse(tokenString string) (*Claims, error) {
	jwtConfig := post_board_config.GetJwtConfig()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtConfig.Secret, nil
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
