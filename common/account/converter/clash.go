package converter

import (
	"fmt"
	"strings"

	"github.com/mangkoyla/LatinaSub-go/db"
	C "github.com/sagernet/sing-box/constant"
)

func ToClash(accounts []db.DBScheme) string {
	result := []string{"proxies:"}
	result = append(result, fmt.Sprintf("    # %s", "Clash have inconsistent configuration for SS between V2ray and OBFS plugin"))
	result = append(result, fmt.Sprintf("    # %s", "So i'm lazy to implement those plugins, unless someone PR it or until i got very motivated"))

	for _, account := range accounts {
		var proxy []string

		proxy = append(proxy, fmt.Sprintf("  - name: %s", account.Remark))
		proxy = append(proxy, fmt.Sprintf("    server: %s", account.Server))
		proxy = append(proxy, fmt.Sprintf("    port: %d", account.ServerPort))

		switch account.VPN {
		case C.TypeVMess, C.TypeVLESS:
			proxy = append(proxy, fmt.Sprintf("    type: %s", account.VPN))
			proxy = append(proxy, fmt.Sprintf("    uuid: %s", account.UUID))
			proxy = append(proxy, fmt.Sprintf("    cipher: %s", "auto"))
			proxy = append(proxy, fmt.Sprintf("    tls: %t", account.TLS))
			proxy = append(proxy, fmt.Sprintf("    udp: %t", true))
			proxy = append(proxy, fmt.Sprintf("    skip-cert-verify: %t", true))
			proxy = append(proxy, fmt.Sprintf("    servername: %s", account.SNI))
			proxy = append(proxy, fmt.Sprintf("    network: %s", account.Transport))

			switch account.VPN {
			case "vmess":
				proxy = append(proxy, fmt.Sprintf("    alterId: %d", account.AlterId))
			}
		case C.TypeTrojan:
			proxy = append(proxy, fmt.Sprintf("    type: %s", account.VPN))
			proxy = append(proxy, fmt.Sprintf("    password: %s", account.Password))
			proxy = append(proxy, fmt.Sprintf("    udp: %t", true))
			proxy = append(proxy, fmt.Sprintf("    skip-cert-verify: %t", true))
			proxy = append(proxy, fmt.Sprintf("    sni: %s", account.SNI))
			proxy = append(proxy, fmt.Sprintf("    network: %s", account.Transport))
		case C.TypeShadowsocks:
			proxy = append(proxy, fmt.Sprintf("    type: %s", "ss"))
			proxy = append(proxy, fmt.Sprintf("    cipher: %s", account.Method))
			proxy = append(proxy, fmt.Sprintf("    password: %s", account.Password))
		case C.TypeShadowsocksR:
			proxy = append(proxy, fmt.Sprintf("    type: %s", "ssr"))
			proxy = append(proxy, fmt.Sprintf("    cipher: %s", account.Method))
			proxy = append(proxy, fmt.Sprintf("    password: %s", account.Password))
			proxy = append(proxy, fmt.Sprintf("    obfs: %s", account.OBFS))
			proxy = append(proxy, fmt.Sprintf("    obfs-param: %s", account.OBFSParam))
			proxy = append(proxy, fmt.Sprintf("    protocol: %s", account.Protocol))
			proxy = append(proxy, fmt.Sprintf("    protocol-param: %s", account.ProtocolParam))
			proxy = append(proxy, fmt.Sprintf("    udp: %t", true))
		}

		switch account.Transport {
		case "ws", "websocket":
			proxy = append(proxy, fmt.Sprintf("    ws-opts: %s", ""))
			proxy = append(proxy, fmt.Sprintf("      path: %s", account.Path))
			proxy = append(proxy, fmt.Sprintf("      headers: %s", ""))
			proxy = append(proxy, fmt.Sprintf("        Host: %s", account.Host))
		case "grpc":
			proxy = append(proxy, fmt.Sprintf("    grpc-opts: %s", ""))
			proxy = append(proxy, fmt.Sprintf("      grpc-service-name: %s", account.ServiceName))
		}

		result = append(result, proxy...)
	}

	return strings.Join(result, "\n")
}
