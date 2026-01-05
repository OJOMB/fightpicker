package jsonwebtokens

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func (jwtt *JWTTool[T]) Parse(tokenString string, key []byte) (*T, *jwt.RegisteredClaims, error) {
	var claims Claims[T]
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return key, nil
	})
	if err != nil {
		return nil, nil, err
	}

	if !token.Valid {
		return nil, nil, errors.New("invalid token")
	}

	return &claims.Custom, &claims.RegisteredClaims, nil
}
