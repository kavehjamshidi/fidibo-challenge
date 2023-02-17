package service

import (
	"log"
	"time"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/internal/token"
)

type RefreshTokenService interface {
	RefreshToken(request domain.RefreshTokenRequest) (domain.RefreshTokenResponse, error)
}

type refreshTokenService struct {
	accessTokenExpiry  time.Duration
	accessTokenSecret  string
	refreshTokenExpiry time.Duration
	refreshTokenSecret string
}

func (l *refreshTokenService) RefreshToken(request domain.RefreshTokenRequest) (domain.RefreshTokenResponse, error) {
	accessToken, err := token.GenerateJWT(request.Username, l.accessTokenSecret, l.accessTokenExpiry)
	if err != nil {
		log.Printf("RefreshToken Service - could not generate access token: %v", err)
		return domain.RefreshTokenResponse{}, err
	}

	refreshToken, err := token.GenerateJWT(request.Username, l.refreshTokenSecret, l.refreshTokenExpiry)
	if err != nil {
		log.Printf("RefreshToken Service - could not generate refresh token: %v", err)
		return domain.RefreshTokenResponse{}, err
	}

	return domain.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewRefreshTokenService(accessTokenExpiry time.Duration,
	accessTokenSecret string,
	refreshTokenExpiry time.Duration,
	refreshTokenSecret string) RefreshTokenService {
	return &refreshTokenService{
		accessTokenExpiry:  accessTokenExpiry,
		accessTokenSecret:  accessTokenSecret,
		refreshTokenExpiry: refreshTokenExpiry,
		refreshTokenSecret: refreshTokenSecret,
	}
}
