package bootstrap

import (
	"log"
	"os"
	"time"
)

const (
	serverAddressEnvKey              = "SERVER_ADDRESS"
	redisURIEnvKey                   = "REDIS_URI"
	accessTokenExpirationTimeEnvKey  = "ACCESS_EXPIRATION"
	refreshTokenExpirationTimeEnvKey = "REFRESH_EXPIRATION"
	accessTokenSecretEnvKey          = "ACCESS_SECRET"
	refreshTokenSecretEnvKey         = "REFRESH_SECRET"

	defaultServerAddress              = ":8080"
	defaultRedisURI                   = "redis://6379"
	defaultAccessTokenExpirationTime  = "15m"
	defaultRefreshTokenExpirationTime = "168h"
	defaultAccessTokenSecret          = "access token secret"
	defaultRefreshTokenSecret         = "refresh token secret"
)

type Env struct {
	ServerAddress              string
	RedisURI                   string
	AccessTokenExpirationTime  time.Duration
	RefreshTokenExpirationTime time.Duration
	AccessTokenSecret          string
	RefreshTokenSecret         string
}

func NewEnv() *Env {
	serverAddress := getEnvWithFallback(serverAddressEnvKey, defaultServerAddress)
	redisURI := getEnvWithFallback(redisURIEnvKey, defaultRedisURI)
	accessTokenSecret := getEnvWithFallback(accessTokenSecretEnvKey, defaultAccessTokenSecret)
	refreshTokenSecret := getEnvWithFallback(refreshTokenSecretEnvKey, defaultRefreshTokenSecret)

	accessTokenExpirationTimeString := getEnvWithFallback(accessTokenExpirationTimeEnvKey, defaultAccessTokenExpirationTime)
	accessTokenExpirationTime, err := time.ParseDuration(accessTokenExpirationTimeString)
	if err != nil {
		log.Fatalf("could not parse access token expiration time: %v\n", err)
	}
	refreshTokenExpirationTimeString := getEnvWithFallback(refreshTokenExpirationTimeEnvKey, defaultRefreshTokenExpirationTime)
	refreshTokenExpirationTime, err := time.ParseDuration(refreshTokenExpirationTimeString)
	if err != nil {
		log.Fatalf("could not parse refresh token expiration time: %v\n", err)
	}

	return &Env{
		ServerAddress:              serverAddress,
		RedisURI:                   redisURI,
		AccessTokenExpirationTime:  accessTokenExpirationTime,
		RefreshTokenExpirationTime: refreshTokenExpirationTime,
		AccessTokenSecret:          accessTokenSecret,
		RefreshTokenSecret:         refreshTokenSecret,
	}
}

func getEnvWithFallback(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
