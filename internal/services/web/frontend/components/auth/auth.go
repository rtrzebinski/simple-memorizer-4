package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

func User(ctx app.Context) (*frontend.User, error) {
	var accessToken string
	ctx.GetState("resp.AccessToken", &accessToken)
	if accessToken == "" {
		return nil, fmt.Errorf("access token is empty")
	}

	// TODO verify token signature with a public key
	// Parse the JWT without verifying the signature
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}

	return &frontend.User{
		ID:    claims["sub"].(string),
		Name:  claims["name"].(string),
		Email: claims["email"].(string),
	}, nil
}
