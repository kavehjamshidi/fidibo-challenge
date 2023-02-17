package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	username := "test username"
	expiry := 10 * time.Minute
	secret := "test secret"

	token, err := GenerateJWT(username, secret, expiry)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestValidateToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		username := "test username"
		expiry := 10 * time.Minute
		secret := "test secret"

		token, err := GenerateJWT(username, secret, expiry)
		assert.NoError(t, err)

		err = ValidateToken(token, secret)
		assert.NoError(t, err)
	})

	t.Run("invalid token", func(t *testing.T) {
		secret := "test secret"
		err := ValidateToken("invalid token", secret)
		assert.Error(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		username := "test username"
		expiry := -10 * time.Minute
		secret := "test secret"

		token, err := GenerateJWT(username, secret, expiry)
		assert.NoError(t, err)

		err = ValidateToken(token, secret)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "token is expired")
	})
}
