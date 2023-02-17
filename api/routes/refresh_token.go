package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
)

const (
	refreshTokenRoute = "/refresh-token"
)

func SetupRefreshTokenRoutes(r *gin.RouterGroup, controller controllers.RefreshTokenController) {
	r.POST(refreshTokenRoute, controller.RefreshToken)
}
