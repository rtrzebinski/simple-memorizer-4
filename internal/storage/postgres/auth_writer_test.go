package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthWriterSuite struct {
	PostgresSuite
	writer *AuthWriter
}

func TestAuthWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(AuthWriterSuite))
}

func (s *AuthWriterSuite) SetupSuite() {
	s.PostgresSuite.SetupSuite()
	s.writer = NewAuthWriter(s.DB)
}

func (s *AuthWriterSuite) TestAuthWriter_StoreUser() {
	ctx := s.T().Context()

	userID, err := s.writer.StoreUser(ctx, "name", "email", "password")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "1", userID)

	user := FetchUserByEmail(s.DB, "email")
	assert.Equal(s.T(), "name", user.Name)
	assert.Equal(s.T(), "email", user.Email)
	assert.Equal(s.T(), "password", user.Password)
}
