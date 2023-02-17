package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/internal/token"
	"github.com/kavehjamshidi/fidibo-challenge/service"
)

type RefreshTokenController interface {
	RefreshToken(c *gin.Context)
}

type refreshTokenController struct {
	secret string
	svc    service.RefreshTokenService
}

func (r *refreshTokenController) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	username, err := token.ExtractUsername(req.RefreshToken, r.secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := r.svc.RefreshToken(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func NewRefreshTokenController(svc service.RefreshTokenService, secret string) RefreshTokenController {
	return &refreshTokenController{
		svc:    svc,
		secret: secret,
	}
}
