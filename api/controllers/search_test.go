package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		query := "test"
		svcMock := &mocks.SearchService{}
		searchController := NewSearchController(svcMock)

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

		expectedJSONResponse, err := json.Marshal(expectedResult)
		assert.NoError(t, err)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header), URL: &url.URL{}}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		q := c.Request.URL.Query()
		q.Add("keyword", query)
		c.Request.URL.RawQuery = q.Encode()

		svcMock.On("Search", c, query).Return(expectedResult, nil)

		searchController.Search(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, string(expectedJSONResponse), string(res))
		svcMock.AssertExpectations(t)
	})

	t.Run("other error", func(t *testing.T) {
		query := "test"
		svcMock := &mocks.SearchService{}
		searchController := NewSearchController(svcMock)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header), URL: &url.URL{}}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		q := c.Request.URL.Query()
		q.Add("keyword", query)
		c.Request.URL.RawQuery = q.Encode()

		errorMsg := "unknown error"

		svcMock.On("Search", c, query).Return(domain.SearchResult{}, errors.New(errorMsg))

		searchController.Search(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.NotEmpty(t, response.Message)
		assert.Contains(t, response.Message, errorMsg)
		svcMock.AssertExpectations(t)
	})

	t.Run("service unavailable", func(t *testing.T) {
		query := "test"
		svcMock := &mocks.SearchService{}
		searchController := NewSearchController(svcMock)

		w := httptest.NewRecorder()

		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request = &http.Request{Header: make(http.Header), URL: &url.URL{}}
		c.Request.Method = http.MethodPost
		c.Request.Header.Set("Content-Type", "application/json")
		q := c.Request.URL.Query()
		q.Add("keyword", query)
		c.Request.URL.RawQuery = q.Encode()

		errorMsg := "service unavailable: some error"

		svcMock.On("Search", c, query).Return(domain.SearchResult{}, errors.New(errorMsg))

		searchController.Search(c)

		res, err := io.ReadAll(w.Body)
		assert.NoError(t, err)

		response := domain.ErrorResponse{}
		err = json.Unmarshal(res, &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusServiceUnavailable, w.Code)
		assert.NotEmpty(t, response.Message)
		assert.Contains(t, response.Message, errorMsg)
		svcMock.AssertExpectations(t)
	})
}
