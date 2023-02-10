package logRoute

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func LogHandler(c *gin.Context) {
	var (
		file, err          = os.ReadFile("log/scrape.log")
		logLines  []string = strings.Split(string(file), "\n")
		logStr    string
	)

	if err != nil {
		logStr = err.Error()
	} else {
		for i := len(logLines) - 1; i > len(logLines)-100; i-- {
			if i <= -1 {
				break
			}

			logStr = logStr + logLines[i] + "\n"
		}
	}

	c.String(http.StatusOK, logStr)
}
