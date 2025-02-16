package auth

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) SignIn(ctx context.Context, email string, password string) (string, string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.Error(2)
}

type WriterMock struct {
	mock.Mock
}

func (m *WriterMock) Register(ctx context.Context, name string, email string, password string) (string, error) {
	args := m.Called(ctx, name, email, password)
	return args.String(0), args.Error(1)
}

func TestService_Register(t *testing.T) {
	readerMock := &ReaderMock{}
	writerMock := &WriterMock{}
	writerMock.On("Register", nil, "name", "email", mock.MatchedBy(func(password string) bool {
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
	readerMock.On("SignIn", nil, "email", "password").Return("name", "userID", nil)

	service := NewService(readerMock, nil)

	accessToken, err := service.SignIn(nil, "email", "password")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
}
