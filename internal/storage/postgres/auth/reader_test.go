package auth

import (
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestReader_SignIn() {
	reader := NewReader()

	name, userID, err := reader.SignIn(nil, "email", "password")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "name", name)
	assert.Equal(s.T(), "userID", userID)
}
