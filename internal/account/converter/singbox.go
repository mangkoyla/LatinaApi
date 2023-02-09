package converter

import (
	"github.com/LalatinaHub/LatinaSub-go/db"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
)

func ToSingBox(accounts []db.DBScheme) []option.Outbound {
	var outbounds []option.Outbound

	for _, account := range accounts {
		outbound := option.Outbound{}
		Server := option.ServerOptions{
			Server:     account.Server,
			ServerPort: uint16(account.ServerPort),
		}
		TLS := option.OutboundTLSOptions{
			Enabled:    account.TLS,
			ServerName: account.SNI,
			Insecure:   true,
			DisableSNI: false,
		}
		Transport := option.V2RayTransportOptions{
			Type: account.Transport,
			WebsocketOptions: option.V2RayWebsocketOptions{
				Path: account.Path,
				Headers: map[string]string{
					"Host": account.Host,
				},
			},
			GRPCOptions: option.V2RayGRPCOptions{
				ServiceName: account.ServiceName,
			},
			QUICOptions: option.V2RayQUICOptions{},
		}

		if account.Transport == "tcp" {
			Transport.Type = ""
		}

		switch account.VPN {
		case C.TypeVMess:
			outbound.VMessOptions = option.VMessOutboundOptions{
				ServerOptions: Server,
				UUID:          account.UUID,
				Security:      account.Security,
				AlterId:       account.AlterId,
				TLS:           &TLS,
				Transport:     &Transport,
			}
		case C.TypeVLESS:
			outbound.VLESSOptions = option.VLESSOutboundOptions{
				ServerOptions: Server,
				UUID:          account.UUID,
				TLS:           &TLS,
				Transport:     &Transport,
			}
		case C.TypeTrojan:
			outbound.TrojanOptions = option.TrojanOutboundOptions{
				ServerOptions: Server,
				Password:      account.Password,
				TLS:           &TLS,
				Transport:     &Transport,
			}
		case C.TypeShadowsocks:
			outbound.ShadowsocksOptions = option.ShadowsocksOutboundOptions{
				ServerOptions: Server,
				Method:        account.Method,
				Password:      account.Password,
				Plugin:        account.Plugin,
				PluginOptions: account.PluginOpts,
			}
		case C.TypeShadowsocksR:
			outbound.ShadowsocksROptions = option.ShadowsocksROutboundOptions{
				ServerOptions: Server,
				Method:        account.Method,
				Password:      account.Password,
				Obfs:          account.OBFS,
				ObfsParam:     account.OBFSParam,
				Protocol:      account.Protocol,
				ProtocolParam: account.ProtocolParam,
			}
		}

		outbounds = append(outbounds, outbound)
	}

	return outbounds
}
