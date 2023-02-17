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

func TestExtractUsername(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		username := "test username"
		expiry := 10 * time.Minute
		secret := "test secret"

		token, err := GenerateJWT(username, secret, expiry)
		assert.NoError(t, err)

		result, err := ExtractUsername(token, secret)
		assert.NoError(t, err)
		assert.Equal(t, username, result)
	})

	t.Run("invalid token", func(t *testing.T) {
		secret := "test secret"
		result, err := ExtractUsername("invalid token", secret)
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("expired token", func(t *testing.T) {
		username := "test username"
		expiry := -10 * time.Minute
		secret := "test secret"

		token, err := GenerateJWT(username, secret, expiry)
		assert.NoError(t, err)

		result, err := ExtractUsername(token, secret)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "token is expired")
		assert.Empty(t, result)
	})
}
