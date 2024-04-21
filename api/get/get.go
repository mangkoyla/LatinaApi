package main

import (
	"net/http"
	"strings"

	apiHelper "github.com/mangkoyla/LatinaApi/api/helper"
	"github.com/mangkoyla/LatinaApi/common/account"
	"github.com/mangkoyla/LatinaApi/common/account/converter"
	"github.com/mangkoyla/LatinaApi/common/helper"
	"github.com/gin-gonic/gin"
)

func Handler(c *gin.Context) {
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
		c.Data(http.StatusOK, "text/plain", []byte(converter.ToClash(proxies)))
	case "surfboard":
		c.Data(http.StatusOK, "text/plain", []byte(strings.Replace(converter.ToSurfboard(proxies), "URL_PLACEHOLDER", apiHelper.GetRequestedURL(c), 1)))
	case "raw":
		c.Data(http.StatusOK, "text/plain", []byte(converter.ToRaw(proxies)))
	default:
		c.JSON(http.StatusOK, proxies)
	}
}
