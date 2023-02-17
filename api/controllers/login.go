package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kavehjamshidi/fidibo-challenge/domain"
	"github.com/kavehjamshidi/fidibo-challenge/service"
)

type LoginController interface {
	Login(c *gin.Context)
}

type loginController struct {
	svc service.LoginService
}

func (l *loginController) Login(c *gin.Context) {
	var req domain.LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := l.svc.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func NewLoginController(svc service.LoginService) LoginController {
	return &loginController{
		svc: svc,
	}
}
