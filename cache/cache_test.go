package cache

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/redis/go-redis/v9"
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
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectSet(key, jsonData, ttl).SetVal(string(jsonData))

		err = cache.Store(context.TODO(), key, val)
		if err != nil {
			t.Error(err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
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
		if err != nil {
			t.Fatal(err)
		}
		errorMsg := "failed to set"

		mock.ExpectSet(key, jsonData, ttl).SetErr(errors.New(errorMsg))

		err = cache.Store(context.TODO(), key, val)
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
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectGet(key).SetVal(string(jsonData))

		cachedVal, err := cache.Get(context.TODO(), key)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(cachedVal, val) {
			t.Errorf("expected cached value to be %v but got %v", val, cachedVal)
		}
	})

	t.Run("key not found", func(t *testing.T) {
		db, mock := redismock.NewClientMock()

		cache := NewCacher(db)

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

		cache := NewCacher(db)

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
