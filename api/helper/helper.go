package apiHelper

import (
	"fmt"
	"net/url"

	"github.com/gin-gonic/gin"
)

func GetRequestedURL(c *gin.Context) string {
	var queries string
	schema := "http://"
	if c.Request.TLS != nil {
		schema = "https://"
	}

	result := schema + c.Request.Host + c.Request.URL.Path
	for key, value := range c.Request.URL.Query() {
		if queries == "" {
			queries = fmt.Sprintf("?%s=%s", key, url.QueryEscape(value[0]))
		} else {
			queries = fmt.Sprintf("%s&%s=%s", queries, key, url.QueryEscape(value[0]))
		}
	}

	return result + queries
}
