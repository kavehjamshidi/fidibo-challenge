package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
)

const (
	loginRoute = "/login"
)

func SetupLoginRoutes(r *gin.RouterGroup, controller controllers.LoginController) {
	r.POST(loginRoute, controller.Login)
}
