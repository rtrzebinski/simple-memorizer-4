package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthReaderSuite struct {
	PostgresSuite
	reader *AuthReader
}

func TestAuthReader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(AuthReaderSuite))
}

func (s *AuthReaderSuite) SetupSuite() {
	s.PostgresSuite.SetupSuite()
	s.reader = NewAuthReader(s.DB)
}

func (s *AuthReaderSuite) TestAuthReader_FetchUser() {
	ctx := s.T().Context()

	CreateUser(s.DB, &User{
		Name:     "name",
		Email:    "email",
		Password: "password",
	})

	name, userID, password, err := s.reader.FetchUser(ctx, "email")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "name", name)
	assert.Equal(s.T(), "1", userID)
	assert.Equal(s.T(), "password", password)
}
