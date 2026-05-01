package auth

import (
	"context"

	"github.com/Nerzal/gocloak/v13"
)

type GoCloak interface {
	LoginClient(ctx context.Context, clientID, clientSecret, realm string, scopes ...string) (*gocloak.JWT, error)
	CreateUser(ctx context.Context, token, realm string, user gocloak.User) (string, error)
	SetPassword(ctx context.Context, token, userID, realm, password string, temporary bool) error
	UpdateUser(ctx context.Context, token, realm string, user gocloak.User) error
	Login(ctx context.Context, clientID, clientSecret, realm, username, password string) (*gocloak.JWT, error)
	RefreshToken(ctx context.Context, refreshToken, clientID, clientSecret, realm string) (*gocloak.JWT, error)
	RevokeToken(ctx context.Context, realm, clientID, clientSecret, refreshToken string) error
}
