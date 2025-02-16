package auth

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	readerMock := &ReaderMock{}
	writerMock := &WriterMock{}
	writerMock.On("StoreUser", nil, "name", "email", mock.MatchedBy(func(password string) bool {
		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(password), []byte("password")))
		return true
	})).Return("userID", nil)

	service := NewService(readerMock, writerMock)

	accessToken, err := service.Register(nil, "name", "email", "password")
	assert.NoError(t, err)

	publicKeyBytes, err := os.ReadFile(keys + "/public.pem")
	assert.NoError(t, err)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	assert.NoError(t, err)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.True(t, token.Valid)
	assert.Equal(t, "userID", claims["sub"])
	assert.Equal(t, "name", claims["name"])
	assert.Equal(t, "email", claims["email"])
}

func TestService_SignIn(t *testing.T) {
	readerMock := &ReaderMock{}
	readerMock.On("FetchUser", nil, "email").Return("name", "userID", "$2a$10$3bAYWOIv0JgCj2xf9hf.beUioHU5jHIYED.hOxKLttWtNWFp7Aq/O", nil)

	service := NewService(readerMock, nil)

	accessToken, err := service.SignIn(nil, "email", "password")
	assert.NoError(t, err)

	publicKeyBytes, err := os.ReadFile(keys + "/public.pem")
	assert.NoError(t, err)
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	assert.NoError(t, err)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.True(t, token.Valid)
	assert.Equal(t, "userID", claims["sub"])
	assert.Equal(t, "name", claims["name"])
	assert.Equal(t, "email", claims["email"])
}
