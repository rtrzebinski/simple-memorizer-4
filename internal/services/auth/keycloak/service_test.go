package keycloak

import (
	"testing"

	"github.com/Nerzal/gocloak/v13"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Register(t *testing.T) {
	mockKC := new(GoCloakMock)
	cfg := Config{Realm: "realm", ClientID: "cid", ClientSecret: "csecret"}
	service := NewService(mockKC, cfg)
	ctx := t.Context()

	adminJWT := &gocloak.JWT{AccessToken: "admin-token"}
	mockKC.On("LoginClient", ctx, cfg.ClientID, cfg.ClientSecret, cfg.Realm).Return(adminJWT, nil)
	mockKC.On("CreateUser", ctx, "admin-token", cfg.Realm, mock.Anything).Return("user-id", nil)
	mockKC.On("SetPassword", ctx, "admin-token", "user-id", cfg.Realm, "pass", false).Return(nil)
	mockKC.On("UpdateUser", ctx, "admin-token", cfg.Realm, mock.Anything).Return(nil)
	jwt := &gocloak.JWT{AccessToken: "at", IDToken: "idt", ExpiresIn: 1, RefreshExpiresIn: 2, RefreshToken: "rt", TokenType: "tt"}
	mockKC.On("Login", ctx, cfg.ClientID, cfg.ClientSecret, cfg.Realm, "email", "pass").Return(jwt, nil)

	tokens, err := service.Register(ctx, "fn", "ln", "email", "pass")
	assert.NoError(t, err)
	assert.Equal(t, Tokens{
		AccessToken:      "at",
		IDToken:          "idt",
		ExpiresIn:        1,
		RefreshExpiresIn: 2,
		RefreshToken:     "rt",
		TokenType:        "tt",
	}, tokens)
}

func TestService_SignIn(t *testing.T) {
	mockKC := new(GoCloakMock)
	cfg := Config{Realm: "realm", ClientID: "cid", ClientSecret: "csecret"}
	service := NewService(mockKC, cfg)
	ctx := t.Context()
	jwt := &gocloak.JWT{AccessToken: "at", IDToken: "idt", ExpiresIn: 1, RefreshExpiresIn: 2, RefreshToken: "rt", TokenType: "tt"}
	mockKC.On("Login", ctx, cfg.ClientID, cfg.ClientSecret, cfg.Realm, "email", "pass").Return(jwt, nil)

	tokens, err := service.SignIn(ctx, "email", "pass")
	assert.NoError(t, err)
	assert.Equal(t, Tokens{
		AccessToken:      "at",
		IDToken:          "idt",
		ExpiresIn:        1,
		RefreshExpiresIn: 2,
		RefreshToken:     "rt",
		TokenType:        "tt",
	}, tokens)
}

func TestService_Refresh(t *testing.T) {
	mockKC := new(GoCloakMock)
	cfg := Config{Realm: "realm", ClientID: "cid", ClientSecret: "csecret"}
	service := NewService(mockKC, cfg)
	ctx := t.Context()
	jwt := &gocloak.JWT{AccessToken: "at", IDToken: "idt", ExpiresIn: 1, RefreshExpiresIn: 2, RefreshToken: "rt", TokenType: "tt"}
	mockKC.On("RefreshToken", ctx, "refresh", cfg.ClientID, cfg.ClientSecret, cfg.Realm).Return(jwt, nil)

	tokens, err := service.Refresh(ctx, "refresh")
	assert.NoError(t, err)
	assert.Equal(t, Tokens{
		AccessToken:      "at",
		IDToken:          "idt",
		ExpiresIn:        1,
		RefreshExpiresIn: 2,
		RefreshToken:     "rt",
		TokenType:        "tt",
	}, tokens)
}

func TestService_Revoke(t *testing.T) {
	mockKC := new(GoCloakMock)
	cfg := Config{Realm: "realm", ClientID: "cid", ClientSecret: "csecret"}
	service := NewService(mockKC, cfg)
	ctx := t.Context()
	mockKC.On("RevokeToken", ctx, cfg.Realm, cfg.ClientID, cfg.ClientSecret, "refresh").Return(nil)

	err := service.Revoke(ctx, "refresh")
	assert.NoError(t, err)
}
