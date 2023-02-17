package fidibosearch

import (
	"context"
	"net/http"
	"testing"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/stretchr/testify/assert"
)

func TestFidiboSearch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		book := domain.Book{
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
		}
		apiRes := fidiboResposne{}
		apiRes.Books.Hits.Hits = append(apiRes.Books.Hits.Hits, struct {
			Source domain.Book `json:"_source"`
		}{
			Source: book,
		})

		route := "/search"

		srv := newMockServer(route, http.StatusOK, apiRes)
		defer srv.Close()

		f := NewFidiboSearcher("q", srv.URL+route)

		res, err := f.Search(context.TODO(), "test query")

		assert.NoError(t, err)
		assert.Len(t, res.Books, 1)
		assert.Equal(t, res.Books[0], book)
	})

	t.Run("fidibo server error", func(t *testing.T) {
		route := "/search"

		srv := newMockServer(route, http.StatusInternalServerError, "")
		defer srv.Close()

		f := NewFidiboSearcher("q", srv.URL+route)

		res, err := f.Search(context.TODO(), "test query")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "did not receive any response")
		assert.Len(t, res.Books, 0)
	})

	t.Run("unmarshall error", func(t *testing.T) {
		route := "/search"

		srv := newMockServer(route, http.StatusOK, "response")
		defer srv.Close()

		f := NewFidiboSearcher("q", srv.URL+route)

		res, err := f.Search(context.TODO(), "test query")

		assert.Error(t, err)
		assert.Len(t, res.Books, 0)
	})
}
