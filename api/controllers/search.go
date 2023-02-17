package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/service"
)

const queryKey = "keyword"

type SearchController interface {
	Search(c *gin.Context)
}

type searchController struct {
	svc service.SearchService
}

func (s *searchController) Search(c *gin.Context) {
	query, _ := c.GetQuery(queryKey)

	res, err := s.svc.Search(c, query)
	if err != nil {
		statusCode := s.mapErrorToStatusCode(err)
		c.JSON(statusCode, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *searchController) mapErrorToStatusCode(err error) int {
	if strings.Contains(err.Error(), "service unavailable") {
		return http.StatusServiceUnavailable
	}
	return http.StatusInternalServerError
}

func NewSearchController(svc service.SearchService) SearchController {
	return &searchController{svc: svc}
}
