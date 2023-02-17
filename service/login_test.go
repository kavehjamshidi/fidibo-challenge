package service

import (
	"testing"
	"time"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	expiry := 10 * time.Minute
	secret := "test secret"

	credentials := domain.LoginRequest{
		Username: "test",
		Password: "test",
	}

	svc := NewLoginService(expiry, secret, expiry, secret)

	result, err := svc.Login(credentials)
	assert.NoError(t, err)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
}
