package parseRoute

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/LalatinaHub/LatinaSub-go/account"
	"github.com/gin-gonic/gin"
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

	if _, err := json.Marshal(accounts); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, accounts)
	}
}
