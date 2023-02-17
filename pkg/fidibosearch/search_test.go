package fidibosearch

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
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

		if err != nil {
			t.Error(err)
		}

		if len(res.Books) != 1 {
			t.Errorf("expected to receive %d books but got %d books", 1, len(res.Books))
		}

		if !reflect.DeepEqual(res.Books[0], book) {
			t.Errorf("expected to receive book %v but got %v", book, res.Books[0])
		}
	})

	t.Run("fidibo server error", func(t *testing.T) {
		route := "/search"

		srv := newMockServer(route, http.StatusInternalServerError, "")
		defer srv.Close()

		f := NewFidiboSearcher("q", srv.URL+route)

		res, err := f.Search(context.TODO(), "test query")

		fmt.Println(res, err)

		if err == nil {
			t.Errorf("expected to receive an error but got no errors")
		}

		if err.Error() != "did not receive any response" {
			t.Errorf("expected error message to be %v but got %v", "did not receive any response", err.Error())
		}

		if len(res.Books) != 0 {
			t.Errorf("expected to receive %d books but got %d books", 0, len(res.Books))
		}
	})

	t.Run("unmarshall error", func(t *testing.T) {
		route := "/search"

		srv := newMockServer(route, http.StatusOK, "response")
		defer srv.Close()

		f := NewFidiboSearcher("q", srv.URL+route)

		res, err := f.Search(context.TODO(), "test query")

		fmt.Println(res, err)

		if err == nil {
			t.Errorf("expected to receive an error but got no errors")
		}

		if len(res.Books) != 0 {
			t.Errorf("expected to receive %d books but got %d books", 0, len(res.Books))
		}
	})
}
