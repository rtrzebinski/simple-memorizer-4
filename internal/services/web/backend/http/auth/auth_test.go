package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserID(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
	})
	signedToken, _ := token.SignedString([]byte("secret")) // podpis nie jest weryfikowany, więc dowolny klucz

	userID, err := UserID(signedToken)

	assert.NoError(t, err)
	assert.Equal(t, "user-123", userID)
}

func TestUserID_InvalidToken(t *testing.T) {
	_, err := UserID("invalid-token")
	assert.Error(t, err)
}
