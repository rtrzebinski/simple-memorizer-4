package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func UserID(accessToken string) (string, error) {
	// TODO verify token signature with a public key
	// Parse the JWT without verifying the signature
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("parse claims")
	}

	return claims["sub"].(string), nil
}
