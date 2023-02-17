package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
)

type NotFoundController interface {
	NotFound(c *gin.Context)
}

type notFoundController struct{}

func (n *notFoundController) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, domain.ErrorResponse{Message: "Not Found"})
}

func NewNotFoundController() NotFoundController {
	return &notFoundController{}
}
