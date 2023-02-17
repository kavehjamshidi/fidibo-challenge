package service

import (
	"context"
	"errors"
	"testing"

	cacheMock "github.com/kavehjamshidi/fidibo-challenge/cache/mocks"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	fidiboMock "github.com/kavehjamshidi/fidibo-challenge/pkg/fidibosearch/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("successful cache hit", func(t *testing.T) {
		query := "test"

		cache := &cacheMock.Cacher{}
		fidiboClient := &fidiboMock.FidiboSearcher{}

		expectedResult := domain.SearchResult{
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

		cache.On("Get", context.TODO(), query).Return(expectedResult, nil)

		svc := NewSearchService(cache, fidiboClient)
		result, err := svc.Search(context.TODO(), query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		fidiboClient.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("cache miss, get response from http client and stored on cache", func(t *testing.T) {
		query := "test"

		cache := &cacheMock.Cacher{}
		fidiboClient := &fidiboMock.FidiboSearcher{}

		expectedResult := domain.SearchResult{
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

		cache.On("Get", context.TODO(), query).Return(domain.SearchResult{}, redis.Nil)
		fidiboClient.On("Search", context.TODO(), query).Return(expectedResult, nil)
		cache.On("Store", context.TODO(), query, expectedResult).Return(nil)

		svc := NewSearchService(cache, fidiboClient)
		result, err := svc.Search(context.TODO(), query)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)

		cache.AssertExpectations(t)
		fidiboClient.AssertExpectations(t)
	})

	t.Run("cache miss, http client error", func(t *testing.T) {
		query := "test"

		cache := &cacheMock.Cacher{}
		fidiboClient := &fidiboMock.FidiboSearcher{}

		expectedResult := domain.SearchResult{}
		expectedError := errors.New("internal server error")

		cache.On("Get", context.TODO(), query).Return(domain.SearchResult{}, redis.Nil)
		fidiboClient.On("Search", context.TODO(), query).Return(domain.SearchResult{}, expectedError)

		svc := NewSearchService(cache, fidiboClient)
		result, err := svc.Search(context.TODO(), query)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, expectedResult, result)
		cache.AssertExpectations(t)
		fidiboClient.AssertExpectations(t)
	})
}
