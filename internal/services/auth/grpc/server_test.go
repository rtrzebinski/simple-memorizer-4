package grpc

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/auth/keycloak"
	"github.com/stretchr/testify/assert"
)

func TestServer_Register(t *testing.T) {
	mockSvc := new(KeycloakServiceMock)
	s := NewServer(mockSvc)
	ctx := t.Context()
	request := &grpc.RegisterRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
	}
	tokens := keycloak.Tokens{
		AccessToken:      "access",
		IDToken:          "id",
		ExpiresIn:        3600,
		RefreshExpiresIn: 7200,
		RefreshToken:     "refresh",
		TokenType:        "Bearer",
	}
	mockSvc.On("Register", ctx, request.FirstName, request.LastName, request.Email, request.Password).Return(tokens, nil)

	resp, err := s.Register(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, tokens.AccessToken, resp.AccessToken)
	assert.Equal(t, tokens.IDToken, resp.IdToken)
	assert.Equal(t, int32(tokens.ExpiresIn), resp.ExpiresIn)
	assert.Equal(t, int32(tokens.RefreshExpiresIn), resp.RefreshExpiresIn)
	assert.Equal(t, tokens.RefreshToken, resp.RefreshToken)
	assert.Equal(t, tokens.TokenType, resp.TokenType)
	mockSvc.AssertExpectations(t)
}

func TestServer_SignIn(t *testing.T) {
	mockSvc := new(KeycloakServiceMock)
	s := NewServer(mockSvc)
	ctx := t.Context()
	request := &grpc.SignInRequest{
		Email:    "john@example.com",
		Password: "password123",
	}
	tokens := keycloak.Tokens{
		AccessToken:      "access",
		IDToken:          "id",
		ExpiresIn:        3600,
		RefreshExpiresIn: 7200,
		RefreshToken:     "refresh",
		TokenType:        "Bearer",
	}
	mockSvc.On("SignIn", ctx, request.Email, request.Password).Return(tokens, nil)

	resp, err := s.SignIn(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, tokens.AccessToken, resp.AccessToken)
	assert.Equal(t, tokens.IDToken, resp.IdToken)
	assert.Equal(t, int32(tokens.ExpiresIn), resp.ExpiresIn)
	assert.Equal(t, int32(tokens.RefreshExpiresIn), resp.RefreshExpiresIn)
	assert.Equal(t, tokens.RefreshToken, resp.RefreshToken)
	assert.Equal(t, tokens.TokenType, resp.TokenType)
	mockSvc.AssertExpectations(t)
}

func TestServer_Refresh(t *testing.T) {
	mockSvc := new(KeycloakServiceMock)
	s := NewServer(mockSvc)
	ctx := t.Context()
	request := &grpc.RefreshRequest{
		RefreshToken: "refresh",
	}
	tokens := keycloak.Tokens{
		AccessToken:      "access",
		IDToken:          "id",
		ExpiresIn:        3600,
		RefreshExpiresIn: 7200,
		RefreshToken:     "refresh",
		TokenType:        "Bearer",
	}
	mockSvc.On("Refresh", ctx, request.RefreshToken).Return(tokens, nil)

	resp, err := s.Refresh(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, tokens.AccessToken, resp.AccessToken)
	assert.Equal(t, tokens.IDToken, resp.IdToken)
	assert.Equal(t, int32(tokens.ExpiresIn), resp.ExpiresIn)
	assert.Equal(t, int32(tokens.RefreshExpiresIn), resp.RefreshExpiresIn)
	assert.Equal(t, tokens.RefreshToken, resp.RefreshToken)
	assert.Equal(t, tokens.TokenType, resp.TokenType)
	mockSvc.AssertExpectations(t)
}

func TestServer_Revoke(t *testing.T) {
	mockSvc := new(KeycloakServiceMock)
	s := NewServer(mockSvc)
	ctx := t.Context()
	request := &grpc.RevokeRequest{
		RefreshToken: "refresh",
	}
	mockSvc.On("Revoke", ctx, request.RefreshToken).Return(nil)

	resp, err := s.Revoke(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockSvc.AssertExpectations(t)
}
