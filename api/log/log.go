package logRoute

import (
	"net/http"

	"github.com/LalatinaHub/LatinaApi/common/helper"
	"github.com/gin-gonic/gin"
)

func LogHandler(c *gin.Context) {
	c.String(http.StatusOK, helper.GetLastLog())
}
