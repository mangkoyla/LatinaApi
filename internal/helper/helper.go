package helper

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func EncodeToBase64(text string) string {
	return strings.TrimSuffix(base64.StdEncoding.EncodeToString([]byte(text)), "=")
}

func BuildFilter(c *gin.Context) string {
	var (
		filter []string
		result string
	)

	for key, value := range c.Request.URL.Query() {
		switch key {
		case "format", "cdn", "sni", "limit": // Ignore special queries
		case "include":
			var includeFilter []string
			includes := strings.Split(value[0], ",")

			for _, include := range includes {
				includeFilter = append(includeFilter, fmt.Sprintf(`REMARK LIKE "%%%s%%"`, include))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(includeFilter[:], " OR ")))
		case "exclude":
			var excludeFilter []string
			excludes := strings.Split(value[0], ",")

			for _, exclude := range excludes {
				excludeFilter = append(excludeFilter, fmt.Sprintf(`REMARK NOT LIKE "%%%s%%"`, exclude))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(excludeFilter[:], " AND ")))
		case "cc":
			var ccFilter []string
			ccs := strings.Split(value[0], ",")

			for _, cc := range ccs {
				ccFilter = append(ccFilter, fmt.Sprintf(`COUNTRY_CODE="%s"`, cc))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(ccFilter[:], " OR ")))
		case "mode":
			var modeFilter []string
			modes := strings.Split(value[0], ",")

			for _, mode := range modes {
				modeFilter = append(modeFilter, fmt.Sprintf(`CONN_MODE LIKE "%%%s%%"`, mode))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(modeFilter[:], " OR ")))
		case "tls":
			tls, _ := strconv.Atoi(value[0])
			filter = append(filter, fmt.Sprintf(`%s=%d`, strings.ToUpper(key), tls))
		case "network", "transport":
			var transportFilter []string
			transports := strings.Split(value[0], ",")

			for _, transport := range transports {
				transportFilter = append(transportFilter, fmt.Sprintf(`TRANSPORT LIKE "%%%s%%"`, transport))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(transportFilter[:], " OR ")))
		default:
			var valueFilter []string
			values := strings.Split(value[0], ",")

			for _, value := range values {
				valueFilter = append(valueFilter, fmt.Sprintf(`%s="%s"`, strings.ToUpper(key), value))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(valueFilter[:], " OR ")))
		}
	}

	result = strings.Join(filter[:], " AND ")
	result = result + " ORDER BY RANDOM()"
	if limit := c.Query("limit"); limit != "" {
		intLimit, _ := strconv.Atoi(limit)
		result = result + fmt.Sprintf(" LIMIT %d", intLimit)
	}

	if result != "" {
		return "WHERE " + result
	} else {
		return result
	}
}
