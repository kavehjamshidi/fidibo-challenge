package service

import (
	"log"
	"time"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/internal/token"
)

type LoginService interface {
	Login(credentials domain.LoginRequest) (domain.LoginResponse, error)
}

type loginService struct {
	accessTokenExpiry  time.Duration
	accessTokenSecret  string
	refreshTokenExpiry time.Duration
	refreshTokenSecret string
}

func (l *loginService) Login(credentials domain.LoginRequest) (domain.LoginResponse, error) {
	accessToken, err := token.GenerateJWT(credentials.Username, l.accessTokenSecret, l.accessTokenExpiry)
	if err != nil {
		log.Printf("Login Service - could not generate access token: %v", err)
		return domain.LoginResponse{}, err
	}

	refreshToken, err := token.GenerateJWT(credentials.Username, l.refreshTokenSecret, l.refreshTokenExpiry)
	if err != nil {
		log.Printf("Login Service - could not generate refresh token: %v", err)
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewLoginService(accessTokenExpiry time.Duration,
	accessTokenSecret string,
	refreshTokenExpiry time.Duration,
	refreshTokenSecret string) LoginService {
	return &loginService{
		accessTokenExpiry:  accessTokenExpiry,
		accessTokenSecret:  accessTokenSecret,
		refreshTokenExpiry: refreshTokenExpiry,
		refreshTokenSecret: refreshTokenSecret,
	}
}
