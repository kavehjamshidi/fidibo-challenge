package bootstrap

import (
	"os"
	"time"
)

const (
	serverAddressEnvKey              = "SERVER_ADDRESS"
	redisAddressEnvKey               = "REDIS_ADDRESS"
	testRedisAddressEnvKey           = "TEST_REDIS_ADDRESS"
	accessTokenExpirationTimeEnvKey  = "ACCESS_EXPIRATION"
	refreshTokenExpirationTimeEnvKey = "REFRESH_EXPIRATION"
	accessTokenSecretEnvKey          = "ACCESS_SECRET"
	refreshTokenSecretEnvKey         = "REFRESH_SECRET"

	defaultServerAddress              = ":8080"
	defaultRedisAddress               = "localhost:6379"
	defaultTestRedisAddress           = "localhost:6379"
	defaultAccessTokenExpirationTime  = "15m"
	defaultRefreshTokenExpirationTime = "168h"
	defaultAccessTokenSecret          = "access token secret"
	defaultRefreshTokenSecret         = "refresh token secret"
)

type Env struct {
	ServerAddress              string
	RedisAddress               string
	TestRedisAddress           string
	AccessTokenExpirationTime  time.Duration
	RefreshTokenExpirationTime time.Duration
	AccessTokenSecret          string
	RefreshTokenSecret         string
}

func NewEnv() *Env {
	serverAddress := getEnvWithFallback(serverAddressEnvKey, defaultServerAddress)
	redisAddress := getEnvWithFallback(redisAddressEnvKey, defaultRedisAddress)
	testRedisAddress := getEnvWithFallback(testRedisAddressEnvKey, defaultTestRedisAddress)
	accessTokenSecret := getEnvWithFallback(accessTokenSecretEnvKey, defaultAccessTokenSecret)
	refreshTokenSecret := getEnvWithFallback(refreshTokenSecretEnvKey, defaultRefreshTokenSecret)

	accessTokenExpirationTimeString := getEnvWithFallback(accessTokenExpirationTimeEnvKey, defaultAccessTokenExpirationTime)
	accessTokenExpirationTime, err := time.ParseDuration(accessTokenExpirationTimeString)
	if err != nil {
		panic(err)
	}
	refreshTokenExpirationTimeString := getEnvWithFallback(refreshTokenExpirationTimeEnvKey, defaultRefreshTokenExpirationTime)
	refreshTokenExpirationTime, err := time.ParseDuration(refreshTokenExpirationTimeString)
	if err != nil {
		panic(err)
	}

	return &Env{
		ServerAddress:              serverAddress,
		RedisAddress:               redisAddress,
		TestRedisAddress:           testRedisAddress,
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
