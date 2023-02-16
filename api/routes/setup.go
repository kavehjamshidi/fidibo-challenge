package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/middleware"
)

func Setup(gin *gin.Engine) {
	publicRouter := gin.Group("")
	SetupLoginRoutes(publicRouter)
	SetupRefreshTokenRoutes(publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.JWTAuth())
	SetupSearchRoutes(protectedRouter)
}
