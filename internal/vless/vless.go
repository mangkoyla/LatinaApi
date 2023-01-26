package vless

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func Get(filter string) []VlessStruct {
	conn := db.Database.Connect()

	query := fmt.Sprintf(`SELECT 
		ADDRESS,
		PORT,
		PASSWORD,
		SECURITY,
		HOST,
		TYPE,
		PATH,
		SERVICE_NAME,
		MODE,
		ALLOW_INSECURE,
		SNI,
		REMARK,
		FLOW,
		LEVEL,
		IS_CDN,
		CC,
		REGION,
		VPN FROM Vless %s;`, filter)
	rows, _ := conn.Query(query)
	defer rows.Close()
	db.Database.Close(conn)

	return toJson(rows)
}

func toJson(rows *sql.Rows) []VlessStruct {
	var result []VlessStruct

	for rows.Next() {
		var address, password, tls, host, network, path, serviceName, mode, sni, remark, cc, region, vpn, flow string
		var port, level int
		var skipCertVerify, cdn bool
		rows.Scan(&address, &port, &password, &tls, &host, &network, &path, &serviceName, &mode, &skipCertVerify, &sni, &remark, &flow, &level, &cdn, &cc, &region, &vpn)

		result = append(result, VlessStruct{
			ADDRESS:          address,
			PORT:             port,
			PASSWORD:         password,
			SECURITY:         tls,
			HOST:             host,
			NETWORK:          network,
			PATH:             path,
			SERVICE_NAME:     serviceName,
			MODE:             mode,
			SKIP_CERT_VERIFY: skipCertVerify,
			SNI:              sni,
			REMARK:           remark,
			FLOW:             flow,
			LEVEL:            level,
			IS_CDN:           cdn,
			CC:               cc,
			REGION:           region,
			VPN:              vpn,
		})
	}

	return result
}

func ToClash(vlesses []VlessStruct) string {
	var result = []string{"proxies:"}

	for _, vless := range vlesses {
		if vless.SECURITY == "tls" {
			vless.TLS = true
		}

		result = append(result, fmt.Sprintf("  - name: %s", vless.REMARK))
		result = append(result, fmt.Sprintf("    server: %s", vless.ADDRESS))
		result = append(result, fmt.Sprintf("    type: %s", vless.VPN))
		result = append(result, fmt.Sprintf("    port: %d", vless.PORT))
		result = append(result, fmt.Sprintf("    uuid: %s", vless.PASSWORD))
		result = append(result, "    cipher: auto")
		result = append(result, "    udp: true")
		result = append(result, fmt.Sprintf("    tls: %t", vless.TLS))
		result = append(result, fmt.Sprintf("    skip-cert-verify: %t", vless.SKIP_CERT_VERIFY))
		result = append(result, fmt.Sprintf("    network: %s", vless.NETWORK))
		result = append(result, fmt.Sprintf("    servername: %s", vless.SNI))
		if vless.NETWORK == "ws" {
			result = append(result, "    ws-opts:")
			result = append(result, fmt.Sprintf("      path: %s", vless.PATH))
			result = append(result, "      headers:")
			result = append(result, fmt.Sprintf("        Host: %s", vless.HOST))
		} else if vless.NETWORK == "grpc" {
			result = append(result, "    grpc-opts:")
			result = append(result, fmt.Sprintf("      grpc-service-name: '%s'", vless.SERVICE_NAME))
		}
	}

	return strings.Join(result[:], "\n")
}

func ToSingBox(vlesses []VlessStruct) string {
	var result []string

	for _, vless := range vlesses {
		var transportObject, tlsObject string

		if vless.SECURITY == "tls" {
			vless.TLS = true
		}

		tlsObject = fmt.Sprintf(`
		{
			"enabled": %t,
			"disable_sni": false,
			"server_name": "%s",
			"insecure": %t
		}`, vless.TLS, vless.SNI, vless.SKIP_CERT_VERIFY)

		if vless.NETWORK == "ws" {
			transportObject = fmt.Sprintf(`
			{
				"type": "ws",
				"path": "%s",
				"headers": {
					"Host": "%s"
				}
			}`, vless.PATH, vless.HOST)
		} else if vless.NETWORK == "grpc" {
			transportObject = fmt.Sprintf(`
			{
				"type": "grpc",
				"service_name": "%s"
			}`, vless.PATH)
		} else {
			transportObject = `{}`
		}

		result = append(result, fmt.Sprintf(`
		{
			"type": "vless",
			"tag": "%s",
			"server": "%s",
			"server_port": %d,
			"uuid": "%s",
			"tls": %s,
			"transport": %s
		}`, vless.REMARK, vless.ADDRESS, vless.PORT, vless.PASSWORD, tlsObject, transportObject))
	}

	return fmt.Sprintf(`
		{
			"outbounds": [%s]
		}`, strings.Join(result[:], ","))
}

func ToRaw(vlesses []VlessStruct) string {
	var result []string

	for _, vless := range vlesses {
		result = append(result, fmt.Sprintf("vless://%s@%s:%d?security=%s&type=%s&host=%s&sni=%s&path=%s&mode=%s&serviceName=%s#%s", vless.PASSWORD, vless.ADDRESS, vless.PORT, vless.SECURITY, vless.NETWORK, vless.HOST, vless.SNI, vless.PATH, vless.MODE, vless.SERVICE_NAME, url.QueryEscape(vless.REMARK)))
	}

	return strings.Join(result[:], "\n")
}
