package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	ttl = 10 * time.Minute
)

type Cache interface {
	Get(context.Context, string) (string, error)
	Store(context.Context, string, string) error
}

type redisCache struct {
	redisClient *redis.Client
}

func (rc *redisCache) Get(ctx context.Context, key string) (string, error) {
	return rc.redisClient.Get(ctx, key).Result()
}

func (rc *redisCache) Store(ctx context.Context, key string, val string) error {
	return rc.redisClient.Set(ctx, key, val, ttl).Err()
}

func NewCache(redisClient *redis.Client) Cache {
	return &redisCache{
		redisClient: redisClient,
	}
}
