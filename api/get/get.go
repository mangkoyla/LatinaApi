package getRoute

import (
	"net/http"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/account"
	"github.com/LalatinaHub/LatinaApi/internal/account/converter"
	"github.com/LalatinaHub/LatinaApi/internal/helper"
	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	format := c.Query("format")
	cdn := strings.Split(c.DefaultQuery("cdn", ""), ",")
	sni := strings.Split(c.DefaultQuery("sni", ""), ",")

	// Build headers and filters
	disposition := "filename=FUCKMETILLDAYLIGHT"
	filter := helper.BuildFilter(c)
	proxies := account.PopulateBugs(account.Get(filter), cdn, sni)

	// Set headers and filters
	c.Header("Content-Disposition", disposition)

	switch format {
	case "clash":
		c.String(http.StatusOK, converter.ToClash(proxies))
	case "singbox":
		c.JSON(http.StatusOK, converter.ToSingBox(proxies))
	case "surfboard":
		c.String(http.StatusOK, converter.ToSurfboard(proxies))
	case "raw":
		c.String(http.StatusOK, converter.ToRaw(proxies))
	default:
		c.JSON(http.StatusOK, proxies)
	}
}
