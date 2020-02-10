package util

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

//
// ParseToken parses a jwt token and returns the claims
//
func ParseToken(tokenString string) (map[string]interface{}, error) {
	token, _ := jwt.Parse(tokenString, nil)
	if token == nil {
		return nil, errors.New("Cannot parse token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Cannot parse token claims")
	}
	return claims, nil
}
