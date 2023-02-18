package bootstrap

import (
	"os"
	"time"
)

const (
	serverAddressEnvKey      = "SERVER_ADDRESS"
	redisAddressEnvKey       = "REDIS_ADDRESS"
	testRedisAddressEnvKey   = "TEST_REDIS_ADDRESS"
	accessTokenExpiryEnvKey  = "ACCESS_EXPIRY"
	refreshTokenExpiryEnvKey = "REFRESH_EXPIRY"
	accessTokenSecretEnvKey  = "ACCESS_SECRET"
	refreshTokenSecretEnvKey = "REFRESH_SECRET"

	defaultServerAddress      = ":8080"
	defaultRedisAddress       = "localhost:6379"
	defaultTestRedisAddress   = "localhost:6379"
	defaultAccessTokenExpiry  = "15m"
	defaultRefreshTokenExpiry = "168h"
	defaultAccessTokenSecret  = "access token secret"
	defaultRefreshTokenSecret = "refresh token secret"
)

type Env struct {
	ServerAddress      string
	RedisAddress       string
	TestRedisAddress   string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	AccessTokenSecret  string
	RefreshTokenSecret string
}

func NewEnv() *Env {
	serverAddress := getEnvWithFallback(serverAddressEnvKey, defaultServerAddress)
	redisAddress := getEnvWithFallback(redisAddressEnvKey, defaultRedisAddress)
	testRedisAddress := getEnvWithFallback(testRedisAddressEnvKey, defaultTestRedisAddress)
	accessTokenSecret := getEnvWithFallback(accessTokenSecretEnvKey, defaultAccessTokenSecret)
	refreshTokenSecret := getEnvWithFallback(refreshTokenSecretEnvKey, defaultRefreshTokenSecret)

	accessTokenExpiryString := getEnvWithFallback(accessTokenExpiryEnvKey, defaultAccessTokenExpiry)
	accessTokenExpiry, err := time.ParseDuration(accessTokenExpiryString)
	if err != nil {
		panic(err)
	}
	refreshTokenExpiryString := getEnvWithFallback(refreshTokenExpiryEnvKey, defaultRefreshTokenExpiry)
	refreshTokenExpiry, err := time.ParseDuration(refreshTokenExpiryString)
	if err != nil {
		panic(err)
	}

	return &Env{
		ServerAddress:      serverAddress,
		RedisAddress:       redisAddress,
		TestRedisAddress:   testRedisAddress,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshTokenExpiry: refreshTokenExpiry,
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
	}
}

func getEnvWithFallback(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return val
}
