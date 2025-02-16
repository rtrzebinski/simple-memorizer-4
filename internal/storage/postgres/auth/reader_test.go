package auth

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestReader_FetchUser() {
	ctx := context.Background()

	reader := NewReader(s.db)

	postgres.CreateUser(s.db, &postgres.User{
		Name:     "name",
		Email:    "email",
		Password: "password",
	})

	name, userID, password, err := reader.FetchUser(ctx, "email")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "name", name)
	assert.Equal(s.T(), "1", userID)
	assert.Equal(s.T(), "password", password)
}
