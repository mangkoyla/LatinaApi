package helper

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func CalculateMode(isCdn, isSni string) (int, int) {
	if isCdn != "" && isSni != "" {
		return 1, 0
	} else if isCdn != "" {
		return 1, 1
	} else if isSni != "" {
		return 0, 0
	} else {
		return 0, 1
	}
}

func BuildFilter(c *gin.Context) string {
	vpn := c.Query("vpn")
	isCdn, isSni := CalculateMode(c.Query("cdn"), c.Query("sni"))

	var (
		filter []string
		result string
		tls    int
	)

	// Default vpn protocol
	if vpn == "" {
		vpn = "vmess"
	}

	for key, value := range c.Request.URL.Query() {
		switch key {
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
		case "region":
			var regionFilter []string
			regions := strings.Split(value[0], ",")

			for _, region := range regions {
				regionFilter = append(regionFilter, fmt.Sprintf(`%s="%s"`, strings.ToUpper(key), region))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(regionFilter[:], " OR ")))
		case "cc":
			var ccFilter []string
			ccs := strings.Split(value[0], ",")

			for _, cc := range ccs {
				ccFilter = append(ccFilter, fmt.Sprintf(`%s="%s"`, strings.ToUpper(key), cc))
			}

			filter = append(filter, fmt.Sprintf("(%s)", strings.Join(ccFilter[:], " OR ")))
		case "tls":
			if value[0] == "1" {
				tls = 1
			} else {
				tls = 0
			}

			if vpn == "vmess" {
				filter = append(filter, fmt.Sprintf(`%s=%d`, strings.ToUpper(key), tls))
			}
		case "network":
			if vpn == "vmess" {
				filter = append(filter, fmt.Sprintf(`%s="%s"`, strings.ToUpper(key), value[0]))
			} else if vpn == "trojan" {
				filter = append(filter, fmt.Sprintf(`TYPE="%s"`, value[0]))
			} else if vpn == "vless" {
				filter = append(filter, fmt.Sprintf(`TYPE="%s"`, value[0]))
			} else if vpn == "ssr" {
				filter = append(filter, fmt.Sprintf(`OBFS LIKE "%%%s%%"`, value[0]))
			}
		}
	}

	result = strings.Join(filter[:], " AND ")

	if vpn != "ssr" {
		if result != "" {
			result = result + fmt.Sprintf(" AND (IS_CDN=%d OR IS_CDN=%d)", isCdn, isSni)
		} else {
			result = fmt.Sprintf("(IS_CDN=%d OR IS_CDN=%d)", isCdn, isSni)
		}
	}

	if result != "" {
		return "WHERE " + result
	} else {
		return result
	}
}
