package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type TokenVerifier struct {
	client *gocloak.GoCloak
	realm  string
}

func NewTokenVerifier(basePath, realm string) *TokenVerifier {
	return &TokenVerifier{
		client: gocloak.NewClient(basePath),
		realm:  realm,
	}
}

func (v *TokenVerifier) VerifyAndUser(ctx context.Context, accessToken string) (*backend.User, error) {
	// Verify token with Keycloak public key (cached)
	token, claims, err := v.client.DecodeAccessToken(ctx, accessToken, v.realm)
	if err != nil {
		return nil, fmt.Errorf("decode access token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	sub, ok := (*claims)["sub"].(string)
	if !ok || sub == "" {
		return nil, errors.New("missing sub claim")
	}

	return &backend.User{
		ID:    (*claims)["sub"].(string),
		Name:  (*claims)["name"].(string),
		Email: (*claims)["email"].(string),
	}, nil
}
