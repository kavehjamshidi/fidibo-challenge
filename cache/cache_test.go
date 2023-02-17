package cache

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	t.Run("successful store", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

		key := "key1"
		val := domain.SearchResult{
			Books: []domain.Book{
				{
					ImageName: "image.jpg",
					Publishers: domain.Publisher{
						Title: "publisher name",
					},
					ID:      "123",
					Title:   "test title",
					Content: "test content",
					Slug:    "test",
					Authors: []domain.Author{
						{Name: "author name"},
					},
				},
			},
		}
		jsonData, err := json.Marshal(val)
		assert.NoError(t, err)

		mock.ExpectSet(key, jsonData, ttl).SetVal(string(jsonData))

		err = cache.Store(context.TODO(), key, val)
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("failed store", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

		key := "key1"
		val := domain.SearchResult{
			Books: []domain.Book{
				{
					ImageName: "image.jpg",
					Publishers: domain.Publisher{
						Title: "publisher name",
					},
					ID:      "123",
					Title:   "test title",
					Content: "test content",
					Slug:    "test",
					Authors: []domain.Author{
						{Name: "author name"},
					},
				},
			},
		}
		jsonData, err := json.Marshal(val)
		assert.NoError(t, err)

		errorMsg := "failed to set"

		mock.ExpectSet(key, jsonData, ttl).SetErr(errors.New(errorMsg))

		err = cache.Store(context.TODO(), key, val)
		assert.Error(t, err)
		assert.ErrorContains(t, err, errorMsg)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func TestGet(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

		key := "key1"
		val := domain.SearchResult{
			Books: []domain.Book{
				{
					ImageName: "image.jpg",
					Publishers: domain.Publisher{
						Title: "publisher name",
					},
					ID:      "123",
					Title:   "test title",
					Content: "test content",
					Slug:    "test",
					Authors: []domain.Author{
						{Name: "author name"},
					},
				},
			},
		}
		jsonData, err := json.Marshal(val)
		assert.NoError(t, err)

		mock.ExpectGet(key).SetVal(string(jsonData))

		cachedVal, err := cache.Get(context.TODO(), key)
		assert.NoError(t, err)
		assert.Equal(t, val, cachedVal)
	})

	t.Run("key not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

		key := "key1"

		mock.ExpectGet(key).RedisNil()

		_, err := cache.Get(context.TODO(), key)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("other redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

		key := "key1"
		errorMsg := "other error"

		mock.ExpectGet(key).SetErr(errors.New(errorMsg))

		_, err := cache.Get(context.TODO(), key)
		assert.Error(t, err)
		assert.ErrorContains(t, err, errorMsg)
	})
}
