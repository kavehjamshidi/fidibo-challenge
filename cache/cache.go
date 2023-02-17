package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/redis/go-redis/v9"
)

const (
	ttl = 10 * time.Minute
)

type Cacher interface {
	Get(context.Context, string) (domain.SearchResult, error)
	Store(context.Context, string, domain.SearchResult) error
}

type redisCache struct {
	redisClient *redis.Client
}

func (rc *redisCache) Get(ctx context.Context, key string) (domain.SearchResult, error) {
	val, err := rc.redisClient.Get(ctx, key).Result()
	if err != nil {
		return domain.SearchResult{}, err
	}

	res := domain.SearchResult{}
	err = json.Unmarshal([]byte(val), &res)
	if err != nil {
		return domain.SearchResult{}, err
	}

	return res, err
}

func (rc *redisCache) Store(ctx context.Context, key string, val domain.SearchResult) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return rc.redisClient.Set(ctx, key, data, ttl).Err()
}

func NewCacher(redisClient *redis.Client) Cacher {
	return &redisCache{
		redisClient: redisClient,
	}
}
