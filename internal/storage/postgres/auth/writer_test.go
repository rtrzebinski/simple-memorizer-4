package auth

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestWriter_StoreUser() {
	ctx := context.Background()

	writer := NewWriter(s.db)

	userID, err := writer.StoreUser(ctx, "name", "email", "password")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "1", userID)

	user := postgres.FetchUserByEmail(s.db, "email")
	assert.Equal(s.T(), "name", user.Name)
	assert.Equal(s.T(), "email", user.Email)
	assert.Equal(s.T(), "password", user.Password)
}
