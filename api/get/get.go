package endpoint

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	apiHelper "github.com/LalatinaHub/LatinaApi/api/helper"
	"github.com/LalatinaHub/LatinaApi/internal/helper"
	"github.com/LalatinaHub/LatinaApi/internal/ssr"
	"github.com/LalatinaHub/LatinaApi/internal/trojan"
	"github.com/LalatinaHub/LatinaApi/internal/vless"
	"github.com/LalatinaHub/LatinaApi/internal/vmess"
	"github.com/gin-gonic/gin"
)

func GetHandler(c *gin.Context) {
	dl := c.Query("dl")
	format := c.DefaultQuery("format", "clash")
	vpn := c.DefaultQuery("vpn", "vmess")
	cdn := strings.Split(c.DefaultQuery("cdn", ""), ",")
	sni := strings.Split(c.DefaultQuery("sni", ""), ",")

	var (
		proxies, resultJson         json.RawMessage
		err                         error
		disposition, filter, result string

		vmesses []vmess.VmessStruct
		vlesses []vless.VlessStruct
		trojans []trojan.TrojanStruct
		ssrs    []ssr.SsrStruct
	)

	// Build headers and filters
	disposition = fmt.Sprintf("filename=fool_%s.txt", format)
	filter = helper.BuildFilter(c)

	if dl == "1" {
		disposition = "attachment; " + disposition
	}

	// Set headers and filters
	c.Header("Content-Disposition", disposition)

	// Prepare for multi protocol support
	vpns := strings.Split(vpn, ",")
	for _, vpnType := range vpns {
		switch vpnType {
		case "vmess":
			proxies, err = json.Marshal(vmess.Get(filter))
			json.Unmarshal(proxies, &vmesses)
			vmesses = vmess.FillBugs(vmesses, cdn, sni)

			if format == "clash" {
				result = vmess.ToClash(vmesses)
			} else if format == "surfboard" {
				result = vmess.ToSurfboard(vmesses)
				result = strings.Replace(result, "URL_PLACEHOLDER", apiHelper.GetRequestedURL(c), 1)
			} else if format == "singbox" {
				json.Unmarshal([]byte(vmess.ToSingBox(vmesses)), &resultJson)
			} else {
				result = vmess.ToRaw(vmesses)
			}
		case "vless":
			proxies, err = json.Marshal(vless.Get(filter))
			json.Unmarshal(proxies, &vlesses)
			vlesses = vless.FillBugs(vlesses, cdn, sni)

			if format == "clash" {
				result = vless.ToClash(vlesses)
			} else if format == "singbox" {
				json.Unmarshal([]byte(vless.ToSingBox(vlesses)), &resultJson)
			} else {
				result = vless.ToRaw(vlesses)
			}
		case "trojan":
			proxies, err = json.Marshal(trojan.Get(filter))
			json.Unmarshal(proxies, &trojans)
			trojans = trojan.FillBugs(trojans, cdn, sni)

			if format == "clash" {
				result = trojan.ToClash(trojans)
			} else if format == "surfboard" {
				result = trojan.ToSurfboard(trojans)
				result = strings.Replace(result, "URL_PLACEHOLDER", apiHelper.GetRequestedURL(c), 1)
			} else if format == "singbox" {
				json.Unmarshal([]byte(trojan.ToSingBox(trojans)), &resultJson)
			} else {
				result = trojan.ToRaw(trojans)
			}
		case "ssr":
			proxies, err = json.Marshal(ssr.Get(filter))
			json.Unmarshal(proxies, &ssrs)
			ssrs = ssr.FillBugs(ssrs, sni)

			if format == "clash" {
				result = ssr.ToClash(ssrs)
			} else if format == "singbox" {
				json.Unmarshal([]byte(ssr.ToSingBox(ssrs)), &resultJson)
			} else {
				result = ssr.ToRaw(ssrs)
			}
		}
	}

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	} else {
		if format == "singbox" {
			c.JSON(http.StatusOK, resultJson)
		} else {
			c.String(http.StatusOK, result)
		}
		return
	}
}
