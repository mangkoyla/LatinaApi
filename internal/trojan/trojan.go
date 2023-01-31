package trojan

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func Get(filter string) []TrojanStruct {
	conn := db.Connect()

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
		METHOD,
		OTA,
		IS_CDN,
		CC,
		REGION,
		VPN FROM Trojan %s;`, filter)
	rows, _ := conn.Query(query)
	defer rows.Close()
	conn.Close()

	return toJson(rows)
}

func toJson(rows *sql.Rows) []TrojanStruct {
	var result []TrojanStruct
	for rows.Next() {
		var address, password, tls, host, network, path, serviceName, mode, sni, remark, cc, region, vpn, method, flow string
		var port, level int
		var skipCertVerify, cdn, ota bool
		rows.Scan(&address, &port, &password, &tls, &host, &network, &path, &serviceName, &mode, &skipCertVerify, &sni, &remark, &flow, &level, &method, &ota, &cdn, &cc, &region, &vpn)

		result = append(result, TrojanStruct{
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
			METHOD:           method,
			OTA:              ota,
			IS_CDN:           cdn,
			CC:               cc,
			REGION:           region,
			VPN:              vpn,
		})
	}

	return result
}

func ToClash(trojans []TrojanStruct) string {
	var result = []string{"proxies:"}
	for _, trojan := range trojans {
		result = append(result, fmt.Sprintf("  - name: %s", trojan.REMARK))
		result = append(result, fmt.Sprintf("    server: %s", trojan.ADDRESS))
		result = append(result, fmt.Sprintf("    type: %s", trojan.VPN))
		result = append(result, fmt.Sprintf("    port: %d", trojan.PORT))
		result = append(result, fmt.Sprintf("    password: %s", trojan.PASSWORD))
		result = append(result, "    udp: true")
		result = append(result, fmt.Sprintf("    skip-cert-verify: %t", trojan.SKIP_CERT_VERIFY))
		result = append(result, fmt.Sprintf("    network: %s", trojan.NETWORK))
		result = append(result, fmt.Sprintf("    sni: %s", trojan.SNI))
		if trojan.NETWORK == "ws" {
			result = append(result, "    ws-opts:")
			result = append(result, fmt.Sprintf("      path: %s", trojan.PATH))
			result = append(result, "      headers:")
			result = append(result, fmt.Sprintf("        Host: %s", trojan.HOST))
		} else if trojan.NETWORK == "grpc" {
			result = append(result, "    grpc-opts:")
			result = append(result, fmt.Sprintf("      grpc-service-name: '%s'", trojan.SERVICE_NAME))
		}
	}

	return strings.Join(result[:], "\n")
}

func ToSurfboard(trojans []TrojanStruct) string {
	var (
		remarks, proxies []string
		result           string
	)
	modes := [4]string{"SELECT", "URL-TEST", "FALLBACK", "LOAD-BALANCE"}

	baseConfig, err := os.ReadFile("resources/config/surfboard.conf")
	if err != nil {
		log.Fatal(err)
	}

	for _, trojan := range trojans {
		remarks = append(remarks, trojan.REMARK)

		if trojan.SECURITY == "tls" {
			trojan.TLS = true
		}

		proxy := fmt.Sprintf("%s=%s,%s,%d,password=%s,udp-relay=true,tls=%t,skip-cert-verify=%t,sni=%s", trojan.REMARK, trojan.VPN, trojan.ADDRESS, trojan.PORT, trojan.PASSWORD, trojan.TLS, trojan.SKIP_CERT_VERIFY, trojan.SNI)

		if trojan.NETWORK == "ws" {
			proxy = fmt.Sprintf("%s,ws=true,ws-path=%s,ws-headers=Host:%s", proxy, trojan.PATH, trojan.HOST)
		}

		proxies = append(proxies, proxy)
	}

	result = strings.Replace(string(baseConfig), "PROXIES_PLACEHOLDER", strings.Join(proxies[:], "\n"), 1)
	for _, mode := range modes {
		result = fmt.Sprintf("%s\n%s", result, fmt.Sprintf("%s=%s,%s", mode, strings.ToLower(mode), strings.Join(remarks[:], ",")))
	}

	return result
}

func ToSingBox(trojans []TrojanStruct) string {
	var result []string

	for _, trojan := range trojans {
		var transportObject, tlsObject string

		if trojan.SECURITY == "tls" {
			trojan.TLS = true
		}

		tlsObject = fmt.Sprintf(`
		{
			"enabled": %t,
			"disable_sni": false,
			"server_name": "%s",
			"insecure": %t
		}`, trojan.TLS, trojan.SNI, trojan.SKIP_CERT_VERIFY)

		if trojan.NETWORK == "ws" {
			transportObject = fmt.Sprintf(`
			{
				"type": "ws",
				"path": "%s",
				"headers": {
					"Host": "%s"
				}
			}`, trojan.PATH, trojan.HOST)
		} else if trojan.NETWORK == "grpc" {
			transportObject = fmt.Sprintf(`
			{
				"type": "grpc",
				"service_name": "%s"
			}`, trojan.SERVICE_NAME)
		} else {
			transportObject = `{}`
		}

		result = append(result, fmt.Sprintf(`
		{
			"type": "trojan",
			"tag": "%s",
			"server": "%s",
			"server_port": %d,
			"password": "%s",
			"tls": %s,
			"transport": %s
		}`, trojan.REMARK, trojan.ADDRESS, trojan.PORT, trojan.PASSWORD, tlsObject, transportObject))
	}

	return fmt.Sprintf(`
		{
			"outbounds": [%s]
		}`, strings.Join(result[:], ","))
}

func ToRaw(trojans []TrojanStruct) string {
	var result []string

	for _, trojan := range trojans {
		result = append(result, fmt.Sprintf("trojan://%s@%s:%d?security=%s&type=%s&host=%s&sni=%s&path=%s&mode=%s&serviceName=%s#%s", trojan.PASSWORD, trojan.ADDRESS, trojan.PORT, trojan.SECURITY, trojan.NETWORK, trojan.HOST, trojan.SNI, trojan.PATH, trojan.MODE, trojan.SERVICE_NAME, url.QueryEscape(trojan.REMARK)))
	}

	return strings.Join(result[:], "\n")
}
