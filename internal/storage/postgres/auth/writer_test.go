package auth

import "github.com/stretchr/testify/assert"

func (s *PostgresSuite) TestWriter_Register() {

	println("TestWriter_Register")

	writer := NewWriter()

	userID, err := writer.Register(nil, "name", "email", "password")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "userID", userID)
}
