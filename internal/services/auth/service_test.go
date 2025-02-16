package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	readerMock := &ReaderMock{}
	writerMock := &WriterMock{}
	writerMock.On("StoreUser", nil, "name", "email", mock.MatchedBy(func(password string) bool {

		println(password)

		assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(password), []byte("password")))
		return true
	})).Return("userID", nil)

	service := NewService(readerMock, writerMock)

	accessToken, err := service.Register(nil, "name", "email", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}

func TestService_SignIn(t *testing.T) {
	readerMock := &ReaderMock{}
	readerMock.On("FetchUser", nil, "email").Return("name", "userID", "$2a$10$3bAYWOIv0JgCj2xf9hf.beUioHU5jHIYED.hOxKLttWtNWFp7Aq/O", nil)

	service := NewService(readerMock, nil)

	accessToken, err := service.SignIn(nil, "email", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}
