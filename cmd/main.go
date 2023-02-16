package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kavehjamshidi/fidibo-challenge/api/routes"
	"github.com/kavehjamshidi/fidibo-challenge/bootstrap"
)

func main() {
	env := bootstrap.NewEnv()

	r := gin.Default()

	routes.Setup(r)

	r.Run(env.ServerAddress)
}
