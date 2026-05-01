package grpc

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth/keycloak"
	"github.com/stretchr/testify/mock"
)

type KeycloakServiceMock struct {
	mock.Mock
}

func (m *KeycloakServiceMock) Register(ctx context.Context, firstName, lastName, email, password string) (keycloak.Tokens, error) {
	args := m.Called(ctx, firstName, lastName, email, password)
	return args.Get(0).(keycloak.Tokens), args.Error(1)
}
func (m *KeycloakServiceMock) SignIn(ctx context.Context, email, password string) (keycloak.Tokens, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(keycloak.Tokens), args.Error(1)
}
func (m *KeycloakServiceMock) Refresh(ctx context.Context, refreshToken string) (keycloak.Tokens, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(keycloak.Tokens), args.Error(1)
}
func (m *KeycloakServiceMock) Revoke(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}
