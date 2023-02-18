package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/controllers"
	"github.com/kavehjamshidi/fidibo-challenge/api/routes"
	"github.com/kavehjamshidi/fidibo-challenge/bootstrap"
	"github.com/kavehjamshidi/fidibo-challenge/cache"
	"github.com/kavehjamshidi/fidibo-challenge/db"
	"github.com/kavehjamshidi/fidibo-challenge/pkg/fidibosearch"
	"github.com/kavehjamshidi/fidibo-challenge/service"
)

const (
	fidiboQueryKey  = "q"
	fidiboSearchURL = "https://search.fidibo.com"
)

func main() {
	env := bootstrap.NewEnv()

	redisClient := db.NewRedisClient(context.Background(), env.RedisAddress)
	cache := cache.NewCacher(redisClient)

	fidiboClient := fidibosearch.NewFidiboSearcher(fidiboQueryKey, fidiboSearchURL)

	loginSVC := service.NewLoginService(env.AccessTokenExpiry,
		env.AccessTokenSecret,
		env.RefreshTokenExpiry,
		env.RefreshTokenSecret)
	refreshTokenSVC := service.NewRefreshTokenService(env.AccessTokenExpiry,
		env.AccessTokenSecret,
		env.RefreshTokenExpiry,
		env.RefreshTokenSecret)
	searchSVC := service.NewSearchService(cache, fidiboClient)

	loginController := controllers.NewLoginController(loginSVC)
	refreshTokenController := controllers.NewRefreshTokenController(refreshTokenSVC, env.RefreshTokenSecret)
	searchController := controllers.NewSearchController(searchSVC)
	notFoundController := controllers.NewNotFoundController()

	r := gin.Default()

	routes.Setup(r, routes.Controllers{
		SearchController:       searchController,
		LoginController:        loginController,
		RefreshTokenController: refreshTokenController,
	}, env.AccessTokenSecret)

	r.NoRoute(notFoundController.NotFound)

	r.Run(env.ServerAddress)
}
