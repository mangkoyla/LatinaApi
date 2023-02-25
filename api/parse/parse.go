package parseRoute

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/gin-gonic/gin"
	C "github.com/sagernet/sing-box/constant"
)

type PostData struct {
	Urls string `form:"urls"`
}

func ParseHandler(c *gin.Context) {
	var (
		data     PostData
		accounts []account.Account
	)

	if err := c.ShouldBind(&data); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	for _, url := range strings.Split(data.Urls, ",") {
		accounts = append(accounts, *account.New(url))
	}

	for _, account := range accounts {
		switch account.Outbound.Type {
		case C.TypeVMess:
			if account.Outbound.VMessOptions.Transport != nil {
				if account.Outbound.VMessOptions.Transport.Type == "" {
					account.Outbound.VMessOptions.Transport = nil
				}
			}
		case C.TypeVLESS:
			if account.Outbound.VLESSOptions.Transport != nil {
				if account.Outbound.VLESSOptions.Transport.Type == "" {
					account.Outbound.VLESSOptions.Transport = nil
				}
			}
		case C.TypeTrojan:
			if account.Outbound.TrojanOptions.Transport != nil {
				if account.Outbound.TrojanOptions.Transport.Type == "" {
					account.Outbound.TrojanOptions.Transport = nil
				}
			}
		}

	}

	if _, err := json.Marshal(accounts); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, accounts)
	}
}
