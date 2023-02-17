package service

import (
	"context"
	"log"

	"github.com/kavehjamshidi/fidibo-challenge/cache"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/pkg/fidibosearch"
)

type SearchService interface {
	Search(ctx context.Context, query string) (domain.SearchResult, error)
}

type searchService struct {
	cache        cache.Cacher
	fidiboSearch fidibosearch.FidiboSearcher
}

func (s *searchService) Search(ctx context.Context, query string) (domain.SearchResult, error) {
	cachedRes, err := s.cache.Get(ctx, query)
	if err == nil {
		return cachedRes, nil
	} else {
		log.Printf("Fidibo Search Cache Retreival Error: %v\n", err)
	}

	fidiboRes, err := s.fidiboSearch.Search(ctx, query)
	if err != nil {
		log.Printf("Fidibo Search Error: %v\n", err)
		return domain.SearchResult{}, err
	}

	err = s.cache.Store(ctx, query, fidiboRes)
	if err != nil {
		log.Printf("Fidibo Search Cache Store Error: %v\n", err)
	}

	return fidiboRes, nil
}

func NewSearchService(cache cache.Cacher, fidiboSearch fidibosearch.FidiboSearcher) SearchService {
	return &searchService{
		cache:        cache,
		fidiboSearch: fidiboSearch,
	}
}
