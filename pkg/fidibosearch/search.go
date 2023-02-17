package fidibosearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
)

type fidiboResposne struct {
	Books struct {
		Hits struct {
			Hits []struct {
				Source domain.Book `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	} `json:"books"`
}

type FidiboSearcher interface {
	Search(context.Context, string) (domain.SearchResult, error)
}

type fidiboSearcher struct {
	queryKey string
	url      string
}

func (f *fidiboSearcher) Search(ctx context.Context, query string) (domain.SearchResult, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.WriteField(f.queryKey, query)
	if err != nil {
		return domain.SearchResult{}, err
	}
	err = writer.Close()
	if err != nil {
		return domain.SearchResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.url, payload)
	if err != nil {
		return domain.SearchResult{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return domain.SearchResult{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.SearchResult{}, errors.New("did not receive any response")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return domain.SearchResult{}, err
	}

	fidiboResponse := fidiboResposne{}
	err = json.Unmarshal(body, &fidiboResponse)
	if err != nil {
		return domain.SearchResult{}, err
	}

	return f.convertFidiboResponseToDomainModel(fidiboResponse), err
}

func NewFidiboSearcher(queryKey, url string) FidiboSearcher {
	return &fidiboSearcher{
		queryKey: queryKey,
		url:      url,
	}
}

func (s *fidiboSearcher) convertFidiboResponseToDomainModel(f fidiboResposne) domain.SearchResult {
	var result domain.SearchResult
	for _, v := range f.Books.Hits.Hits {
		result.Books = append(result.Books, v.Source)
	}
	return result
}
