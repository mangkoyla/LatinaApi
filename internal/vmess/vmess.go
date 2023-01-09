package vmess

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func GetAll(isCdn, isSni int) []VmessStruct {
	conn := db.Database.Connect()

	query := fmt.Sprintf(`SELECT 
		ADDRESS,
		ALTER_ID,
		PORT,
		PASSWORD,
		SECURITY,
		HOST,
		TLS,
		NETWORK,
		PATH,
		SKIP_CERT_VERIFY,
		SNI,
		REMARK,
		IS_CDN,
		CC,
		REGION,
		VPN FROM Vmess WHERE IS_CDN=%d OR IS_CDN=%d;`, isCdn, isSni)
	rows, _ := conn.Query(query)
	defer rows.Close()
	db.Database.Close(conn)

	return toJson(rows)
}

func GetByCC(cc string, isCdn, isSni int) []VmessStruct {
	conn := db.Database.Connect()

	query := fmt.Sprintf(`SELECT 
		ADDRESS,
		ALTER_ID,
		PORT,
		PASSWORD,
		SECURITY,
		HOST,
		TLS,
		NETWORK,
		PATH,
		SKIP_CERT_VERIFY,
		SNI,
		REMARK,
		IS_CDN,
		CC,
		REGION,
		VPN FROM Vmess WHERE CC='%s' AND (IS_CDN=%d OR IS_CDN=%d);`, cc, isCdn, isSni)
	rows, _ := conn.Query(query)
	defer rows.Close()
	db.Database.Close(conn)

	return toJson(rows)
}

func GetByRegion(region string, isCdn, isSni int) []VmessStruct {
	conn := db.Database.Connect()

	query := fmt.Sprintf(`SELECT 
		ADDRESS,
		ALTER_ID,
		PORT,
		PASSWORD,
		SECURITY,
		HOST,
		TLS,
		NETWORK,
		PATH,
		SKIP_CERT_VERIFY,
		SNI,
		REMARK,
		IS_CDN,
		CC,
		REGION,
		VPN FROM Vmess WHERE REGION='%s' AND (IS_CDN=%d OR IS_CDN=%d);`, region, isCdn, isSni)
	rows, _ := conn.Query(query)
	defer rows.Close()
	db.Database.Close(conn)

	return toJson(rows)
}

func toJson(rows *sql.Rows) []VmessStruct {
	var result []VmessStruct
	for rows.Next() {
		var address, password, security, host, network, path, sni, remark, cc, region, vpn string
		var alterId, port int
		var tls, skipCertVerify, cdn bool
		rows.Scan(&address, &alterId, &port, &password, &security, &host, &tls, &network, &path, &skipCertVerify, &sni, &remark, &cdn, &cc, &region, &vpn)

		result = append(result, VmessStruct{
			ADDRESS:          address,
			ALTER_ID:         alterId,
			PORT:             port,
			PASSWORD:         password,
			SECURITY:         security,
			HOST:             host,
			TLS:              tls,
			NETWORK:          network,
			PATH:             path,
			SKIP_CERT_VERIFY: skipCertVerify,
			SNI:              sni,
			REMARK:           remark,
			IS_CDN:           cdn,
			CC:               cc,
			REGION:           region,
			VPN:              vpn,
		})
	}

	return result
}

func ToClash(vmesses []VmessStruct) string {
	var result = []string{"proxies:"}
	for _, vmess := range vmesses {
		result = append(result, fmt.Sprintf("  - name: %s", vmess.REMARK))
		result = append(result, fmt.Sprintf("    server: %s", vmess.ADDRESS))
		result = append(result, fmt.Sprintf("    type: %s", vmess.VPN))
		result = append(result, fmt.Sprintf("    port: %d", vmess.PORT))
		result = append(result, fmt.Sprintf("    uuid: %s", vmess.PASSWORD))
		result = append(result, fmt.Sprintf("    alterId: %d", vmess.ALTER_ID))
		result = append(result, "    cipher: auto")
		result = append(result, fmt.Sprintf("    tls: %t", vmess.TLS))
		result = append(result, "    udp: true")
		result = append(result, fmt.Sprintf("    skip-cert-verify: %t", vmess.SKIP_CERT_VERIFY))
		result = append(result, fmt.Sprintf("    network: %s", vmess.NETWORK))
		result = append(result, fmt.Sprintf("    servername: %s", vmess.SNI))
		if vmess.NETWORK == "ws" {
			result = append(result, "    ws-opts:")
			result = append(result, fmt.Sprintf("      path: %s", vmess.PATH))
			result = append(result, "      headers:")
			result = append(result, fmt.Sprintf("        Host: %s", vmess.HOST))
		} else if vmess.NETWORK == "grpc" {
			result = append(result, "    grpc-opts:")
			result = append(result, fmt.Sprintf("      grpc-service-name: '%s'", vmess.PATH))
		}
	}

	return strings.Join(result[:], "\n")
}

func ToSurfboard(vmesses []VmessStruct) string {
	var remarks []string
	var proxies []string
	var result string
	modes := [4]string{"SELECT", "URL-TEST", "FALLBACK", "LOAD-BALANCE"}

	baseConfig, err := os.ReadFile("resources/config/surfboard.conf")
	if err != nil {
		log.Fatal(err)
	}

	for _, vmess := range vmesses {
		remarks = append(remarks, vmess.REMARK)
		proxy := fmt.Sprintf("%s=%s,%s,%d,username=%s,udp-relay=true,tls=%t,skip-cert-verify=%t,sni=%s", vmess.REMARK, vmess.VPN, vmess.ADDRESS, vmess.PORT, vmess.PASSWORD, vmess.TLS, vmess.SKIP_CERT_VERIFY, vmess.SNI)

		if vmess.NETWORK == "ws" {
			proxy = fmt.Sprintf("%s,ws=true,ws-path=%s,ws-headers=Host:%s", proxy, vmess.PATH, vmess.HOST)
		}

		proxies = append(proxies, proxy)
	}

	result = strings.Replace(string(baseConfig), "PROXIES_PLACEHOLDER", strings.Join(proxies[:], "\n"), 1)
	for _, mode := range modes {
		result = fmt.Sprintf("%s\n%s", result, fmt.Sprintf("%s=%s,%s", mode, strings.ToLower(mode), strings.Join(remarks[:], ",")))
	}

	return result
}

func ToRaw(vmesses []VmessStruct) string {
	var result []string

	for _, vmess := range vmesses {
		var tls string = "tls"
		if !vmess.TLS {
			tls = ""
		}

		var vmessJson json.RawMessage
		vmessJsonByte := []byte(fmt.Sprintf(`{"add": "%s","aid": %d,"host": "%s","id": "%s","net": "%s","path": "%s","port": %d,"ps": "%s","tls": "%s","security": "%s","skip-cert-verify": %t,"sni": "%s"}`, vmess.ADDRESS, vmess.ALTER_ID, vmess.HOST, vmess.PASSWORD, vmess.NETWORK, vmess.PATH, vmess.PORT, vmess.REMARK, tls, vmess.SECURITY, vmess.SKIP_CERT_VERIFY, vmess.SNI))

		err := json.Unmarshal(vmessJsonByte, &vmessJson)
		if err != nil {
			log.Print(err)
		}

		result = append(result, "vmess://"+base64.StdEncoding.EncodeToString(vmessJson))
	}

	return strings.Join(result[:], "\n")
}
