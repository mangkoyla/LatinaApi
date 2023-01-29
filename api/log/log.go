package logRoute

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func LogHandler(c *gin.Context) {
	logText, _ := os.ReadFile("LatinaSub/log.txt")

	c.String(http.StatusOK, string(logText))
}
