package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
)

func TestStore(t *testing.T) {
	t.Run("successful store", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCache(db)

		key := "key1"
		val := "val1"

		mock.ExpectSet(key, val, ttl).SetVal(val)

		err := cache.Store(context.TODO(), key, val)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})

	t.Run("failed store", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCache(db)

		key := "key1"
		val := "val1"
		errorMsg := "failed to set"

		mock.ExpectSet(key, val, ttl).SetErr(errors.New(errorMsg))

		err := cache.Store(context.TODO(), key, val)
		if err == nil {
			t.Error("expected error but got no errors")
		}

		if err.Error() != errorMsg {
			t.Errorf("expected error message to be %v but got %v", errorMsg, err.Error())
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCache(db)

		key := "key1"
		val := "val1"

		mock.ExpectGet(key).SetVal(val)

		cachedVal, err := cache.Get(context.TODO(), key)
		if err != nil {
			t.Error(err)
		}

		if cachedVal != val {
			t.Errorf("expected cached value to be %s but got %s", val, cachedVal)
		}
	})

	t.Run("key not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCache(db)

		key := "key1"

		mock.ExpectGet(key).RedisNil()

		_, err := cache.Get(context.TODO(), key)
		if err == nil {
			t.Error("expected RedisNil error but got no errors")
		}

		if err != redis.Nil {
			t.Errorf("expected error to be %v but got %v", redis.Nil, err)
		}
	})

	t.Run("other redis error", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCache(db)

		key := "key1"
		errorMsg := "other error"

		mock.ExpectGet(key).SetErr(errors.New(errorMsg))

		_, err := cache.Get(context.TODO(), key)
		if err == nil {
			t.Error("expected error but got no errors")
		}

		if err.Error() != errorMsg {
			t.Errorf("expected error message to be %v but got %v", errorMsg, err.Error())
		}
	})
}
