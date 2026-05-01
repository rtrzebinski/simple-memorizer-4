package grpc

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth/keycloak"
)

type KeycloakService interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (keycloak.Tokens, error)
	SignIn(ctx context.Context, email, password string) (keycloak.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (keycloak.Tokens, error)
	Revoke(ctx context.Context, refreshToken string) error
}
