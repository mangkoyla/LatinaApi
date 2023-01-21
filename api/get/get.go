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
	cc := c.DefaultQuery("cc", "all")
	region := c.DefaultQuery("region", "all")
	cdn := strings.Split(c.DefaultQuery("cdn", ""), ",")
	sni := strings.Split(c.DefaultQuery("sni", ""), ",")
	isCdn, isSni := helper.CalculateMode(c.Query("cdn"), c.Query("sni"))

	var (
		proxies     json.RawMessage
		err         error
		disposition string
	)

	if cc != "all" {
		disposition = "filename=" + strings.ToUpper(fmt.Sprintf("%s_%s_%s", format, cc, vpn))
	} else if region != "all" {
		disposition = "filename=" + strings.ToUpper(fmt.Sprintf("%s_%s_%s", format, region, vpn))
	} else {
		disposition = "filename=" + strings.ToUpper(fmt.Sprintf("%s_all_%s", format, vpn))
	}

	if dl != "" {
		disposition = "attachment; " + disposition
	}

	// Set headers
	c.Header("Content-Disposition", disposition)

	if vpn == "vmess" {
		if cc != "all" {
			proxies, err = json.Marshal(vmess.GetByCC(cc, isCdn, isSni))
		} else if region != "all" {
			proxies, err = json.Marshal(vmess.GetByRegion(region, isCdn, isSni))
		} else {
			proxies, err = json.Marshal(vmess.GetAll(isCdn, isSni))
		}

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var vmesses []vmess.VmessStruct
		json.Unmarshal(proxies, &vmesses)

		vmesses = vmess.FillBugs(vmesses, cdn, sni)
		if format == "clash" {
			result := vmess.ToClash(vmesses)
			c.String(http.StatusOK, result)
		} else if format == "surfboard" {
			result := vmess.ToSurfboard(vmesses)
			result = strings.Replace(result, "URL_PLACEHOLDER", apiHelper.GetRequestedURL(c), 1)
			c.String(http.StatusOK, result)
		} else if format == "singbox" {
			var result json.RawMessage
			json.Unmarshal([]byte(vmess.ToSingBox(vmesses)), &result)
			c.JSON(http.StatusOK, result)
		} else if format == "raw" {
			result := vmess.ToRaw(vmesses)
			c.String(http.StatusOK, result)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		return
	} else if vpn == "trojan" {
		if cc != "all" {
			proxies, err = json.Marshal(trojan.GetByCC(cc, isCdn, isSni))
		} else if region != "all" {
			proxies, err = json.Marshal(trojan.GetByRegion(region, isCdn, isSni))
		} else {
			proxies, err = json.Marshal(trojan.GetAll(isCdn, isSni))
		}

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var trojans []trojan.TrojanStruct
		json.Unmarshal(proxies, &trojans)

		trojans = trojan.FillBugs(trojans, cdn, sni)
		if format == "clash" {
			result := trojan.ToClash(trojans)
			c.String(http.StatusOK, result)
		} else if format == "surfboard" {
			result := trojan.ToSurfboard(trojans)
			result = strings.Replace(result, "URL_PLACEHOLDER", apiHelper.GetRequestedURL(c), 1)
			c.String(http.StatusOK, result)
		} else if format == "singbox" {
			var result json.RawMessage
			json.Unmarshal([]byte(trojan.ToSingBox(trojans)), &result)
			c.JSON(http.StatusOK, result)
		} else if format == "raw" {
			result := trojan.ToRaw(trojans)
			c.String(http.StatusOK, result)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		return
	} else if vpn == "ssr" {
		if cc != "all" {
			proxies, err = json.Marshal(ssr.GetByCC(cc))
		} else if region != "all" {
			proxies, err = json.Marshal(ssr.GetByRegion(region))
		} else {
			proxies, err = json.Marshal(ssr.GetAll())
		}

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var ssrs []ssr.SsrStruct
		json.Unmarshal(proxies, &ssrs)

		ssrs = ssr.FillBugs(ssrs, sni)
		if format == "clash" {
			result := ssr.ToClash(ssrs)
			c.String(http.StatusOK, result)
		} else if format == "singbox" {
			var result json.RawMessage
			json.Unmarshal([]byte(ssr.ToSingBox(ssrs)), &result)
			c.JSON(http.StatusOK, result)
		} else if format == "raw" {
			result := ssr.ToRaw(ssrs)
			c.String(http.StatusOK, result)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		return
	} else if vpn == "vless" {
		if cc != "all" {
			proxies, err = json.Marshal(vless.GetByCC(cc, isCdn, isSni))
		} else if region != "all" {
			proxies, err = json.Marshal(vless.GetByRegion(region, isCdn, isSni))
		} else {
			proxies, err = json.Marshal(vless.GetAll(isCdn, isSni))
		}

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		var vlesses []vless.VlessStruct
		json.Unmarshal(proxies, &vlesses)

		vlesses = vless.FillBugs(vlesses, cdn, sni)
		if format == "clash" {
			result := vless.ToClash(vlesses)
			c.String(http.StatusOK, result)
		} else if format == "singbox" {
			var result json.RawMessage
			json.Unmarshal([]byte(vless.ToSingBox(vlesses)), &result)
			c.JSON(http.StatusOK, result)
		} else if format == "raw" {
			result := vless.ToRaw(vlesses)
			c.String(http.StatusOK, result)
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		return
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
}
