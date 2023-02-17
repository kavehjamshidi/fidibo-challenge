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
	Search(ctx context.Context, query string) (domain.SearchResult, error)
}

type fidiboClient struct {
	queryKey string
	url      string
}

func (f *fidiboClient) Search(ctx context.Context, query string) (domain.SearchResult, error) {
	req, err := f.createHTTPRequest(ctx, query)
	if err != nil {
		return domain.SearchResult{}, err
	}

	res, err := f.doHTTPRequest(req)
	if err != nil {
		return domain.SearchResult{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return domain.SearchResult{}, errors.New("did not receive any response")
	}

	return f.parseResponse(res.Body)
}

func (f *fidiboClient) convertFidiboResponseToDomainModel(r fidiboResposne) domain.SearchResult {
	var result domain.SearchResult
	for _, v := range r.Books.Hits.Hits {
		result.Books = append(result.Books, v.Source)
	}
	return result
}

func (f *fidiboClient) createHTTPRequest(ctx context.Context, query string) (*http.Request, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	if query != "" {
		err := writer.WriteField(f.queryKey, query)
		if err != nil {
			return nil, err
		}
		err = writer.Close()
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

func (f *fidiboClient) doHTTPRequest(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	return client.Do(req)
}

func (f *fidiboClient) parseResponse(body io.Reader) (domain.SearchResult, error) {
	res, err := io.ReadAll(body)
	if err != nil {
		return domain.SearchResult{}, err
	}

	fidiboResponse := fidiboResposne{}
	err = json.Unmarshal(res, &fidiboResponse)
	if err != nil {
		return domain.SearchResult{}, err
	}

	return f.convertFidiboResponseToDomainModel(fidiboResponse), nil
}

func NewFidiboSearcher(queryKey, url string) FidiboSearcher {
	return &fidiboClient{
		queryKey: queryKey,
		url:      url,
	}
}
