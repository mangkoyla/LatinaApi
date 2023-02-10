package converter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/account/protocol"
	"github.com/LalatinaHub/LatinaApi/internal/helper"
	"github.com/LalatinaHub/LatinaSub-go/db"
	C "github.com/sagernet/sing-box/constant"
)

func ToRaw(accounts []db.DBScheme) string {
	var result []string

	for _, account := range accounts {
		tls := ""
		if account.TLS {
			tls = "tls"
		}

		switch account.VPN {
		case C.TypeVMess:
			tls := ""
			if account.TLS {
				tls = "tls"
			}

			vmess := &protocol.VmessStruct{
				Add:            account.Server,
				Port:           uint16(account.ServerPort),
				Aid:            account.AlterId,
				Id:             account.UUID,
				Net:            account.Transport,
				Path:           account.Path,
				Ps:             account.Remark,
				Tls:            tls,
				Security:       account.Security,
				SkipCretVerify: account.Insecure,
				Sni:            account.SNI,
				Host:           account.Host,
			}
			j, _ := json.Marshal(vmess)
			result = append(result, "vmess://"+helper.EncodeToBase64(string(j)))
		case C.TypeVLESS:
			result = append(result, fmt.Sprintf("vless://%s@%s:%d?security=%s&type=%s&host=%s&sni=%s&path=%s&serviceName=%s#%s", account.UUID, account.Server, account.ServerPort, tls, account.Transport, account.Host, account.SNI, account.Path, account.ServiceName, url.QueryEscape(account.Remark)))
		case C.TypeTrojan:
			result = append(result, fmt.Sprintf("trojan://%s@%s:%d?security=%s&type=%s&host=%s&sni=%s&path=%s&serviceName=%s#%s", account.Password, account.Server, account.ServerPort, tls, account.Transport, account.Host, account.SNI, account.Path, account.ServiceName, url.QueryEscape(account.Remark)))
		case C.TypeShadowsocks:
			result = append(result, "ss://"+helper.EncodeToBase64(fmt.Sprintf("%s:%s@%s:%d", account.Method, account.Password, account.Server, account.ServerPort))+fmt.Sprintf("/?plugin=%s#%s", account.Plugin+";"+account.PluginOpts, url.QueryEscape(account.Remark)))
		case C.TypeShadowsocksR:
			password := helper.EncodeToBase64(account.Password)
			remarks := helper.EncodeToBase64(account.Remark)
			protoParam := helper.EncodeToBase64(account.ProtocolParam)
			obfsParam := helper.EncodeToBase64(account.OBFSParam)

			result = append(result, "ssr://"+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d:%s:%s:%s:%s/?remarks=%s&protoparam=%s&obfsparam=%s", account.Server, account.ServerPort, account.Protocol, account.Method, account.OBFS, password, remarks, protoParam, obfsParam))))
		}
	}

	return strings.Join(result, "\n")
}
