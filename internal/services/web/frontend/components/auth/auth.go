package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

func User(ctx app.Context) (*frontend.User, error) {
	accessToken, err := Token(ctx)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}

	// TODO verify token signature with a public key
	// Parse the JWT without verifying the signature
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("parse claims")
	}

	return &frontend.User{
		ID:    claims["sub"].(string),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
	}, nil
}

func Token(ctx app.Context) (string, error) {
	var accessToken string
	ctx.GetState("AccessToken", &accessToken)
	if accessToken == "" {
		return "", fmt.Errorf("access token is empty")
	}

	return accessToken, nil
}
