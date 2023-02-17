package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
	"github.com/kavehjamshidi/fidibo-challenge/api/middleware"
)

type Controllers struct {
	controllers.SearchController
	controllers.LoginController
	controllers.RefreshTokenController
}

func Setup(gin *gin.Engine, ctrl Controllers, accessTokenSecret string) {
	publicRouter := gin.Group("")
	SetupLoginRoutes(publicRouter, ctrl.LoginController)
	SetupRefreshTokenRoutes(publicRouter, ctrl.RefreshTokenController)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JWTAuth(accessTokenSecret))
	SetupSearchRoutes(protectedRouter, ctrl.SearchController)
}
