package router

import (
	"fmt"
	"net/http"
	"os"

	endpoint "github.com/LalatinaHub/LatinaApi/api/get"
	"github.com/gin-gonic/gin"
)

var (
	Router = gin.Default()
	Port   = os.Getenv("PORT")
)

func Start() {
	Router.SetTrustedProxies(nil)

	Router.StaticFile("/favicon.ico", "resources/favicon.ico")

	Router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": "Go to /get to get nodes",
			"query": gin.H{
				"cc":     "Country Code. ex: SG, ID, JP (select one)",
				"region": "Regional. ex: Asia, Americas, Europe, Africa (select one)",
				"vpn":    "Vpn Type, ex: trojan, vmess, vless (select one)",
				"format": "Target format, ex: clash, surfboard, raw (select one)",
				"cdn":    "CDN bugs, separate with (,) for multiple bugs",
				"sni":    "SNI bugs, separate with (,) for multiple bugs",
			},
			"info":    "Parameter is case sensitive",
			"example": "http://fool.azurewebsites.net/get?vpn=trojan&cdn=hohm.microsoft.com&sni=google.com&cc=SG&format=raw",
		})
	})

	Router.GET("/get", endpoint.GetHandler)

	if Port == "" {
		Port = "8080"
	}
	fmt.Println("Server listening on port " + Port)
	Router.Run(":" + Port)
}
