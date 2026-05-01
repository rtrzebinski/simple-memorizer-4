package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/mock"
)

type GoCloakMock struct {
	mock.Mock
}

func (m *GoCloakMock) LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*gocloak.JWT, error) {
	args := m.Called(ctx, clientID, clientSecret, realm)
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}
func (m *GoCloakMock) CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error) {
	args := m.Called(ctx, token, realm, user)
	return args.String(0), args.Error(1)
}
func (m *GoCloakMock) SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error {
	args := m.Called(ctx, token, userID, realm, password, temporary)
	return args.Error(0)
}
func (m *GoCloakMock) UpdateUser(ctx context.Context, token, realm string, user gocloak.User) error {
	args := m.Called(ctx, token, realm, user)
	return args.Error(0)
}
func (m *GoCloakMock) Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*gocloak.JWT, error) {
	args := m.Called(ctx, clientID, clientSecret, realm, username, password)
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}
func (m *GoCloakMock) RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error) {
	args := m.Called(ctx, refreshToken, clientID, clientSecret, realm)
	return args.Get(0).(*gocloak.JWT), args.Error(1)
}
func (m *GoCloakMock) RevokeToken(ctx context.Context, realm, clientID, clientSecret, refreshToken string) error {
	args := m.Called(ctx, realm, clientID, clientSecret, refreshToken)
	return args.Error(0)
}
