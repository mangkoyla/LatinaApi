package router

import (
	"fmt"
	"os"

	endpoint "github.com/LalatinaHub/LatinaApi/api/get"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	Router = gin.Default()
	Port   = os.Getenv("PORT")
)

func Start() {
	Router.SetTrustedProxies(nil)

	Router.Use(static.Serve("/", static.LocalFile("public/", false)))

	Router.GET("/get", endpoint.GetHandler)

	if Port == "" {
		Port = "8080"
	}
	fmt.Println("Server listening on port " + Port)
	Router.Run(":" + Port)
}
