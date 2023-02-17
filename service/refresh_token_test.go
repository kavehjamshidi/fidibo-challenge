package service

import (
	"testing"
	"time"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/stretchr/testify/assert"
)

func TestRefreshToken(t *testing.T) {
	expiry := 10 * time.Minute
	secret := "test secret"
	username := "test"

	svc := NewRefreshTokenService(expiry, secret, expiry, secret)

	result, err := svc.RefreshToken(domain.RefreshTokenRequest{Username: username})
	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}
