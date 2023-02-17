package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
)

const (
	searchRoute = "/search/book"
)

func SetupSearchRoutes(r *gin.RouterGroup, controller controllers.SearchController) {
	r.POST(searchRoute, controller.Search)
}
