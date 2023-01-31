package ssr

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/LalatinaHub/LatinaApi/internal/db"
)

func Get(filter string) []SsrStruct {
	conn := db.Connect()

	query := fmt.Sprintf(`SELECT 
		ADDRESS,
		PORT,
		PASSWORD,
		METHOD,
		PROTOCOL,
		PROTOCOL_PARAM,
		OBFS,
		OBFS_PARAM,
		[GROUP],
		REMARK,
		CC,
		REGION,
		VPN FROM SSR %s`, filter)
	rows, _ := conn.Query(query)
	defer rows.Close()
	conn.Close()

	return toJson(rows)
}

func toJson(rows *sql.Rows) []SsrStruct {
	var result []SsrStruct
	for rows.Next() {
		var address, method, protocol, protocolParam, obfs, obfsParam, group, password, remark, cc, region, vpn string
		var port int
		rows.Scan(&address, &port, &password, &method, &protocol, &protocolParam, &obfs, &obfsParam, &group, &remark, &cc, &region, &vpn)

		result = append(result, SsrStruct{
			ADDRESS:        address,
			PORT:           port,
			PASSWORD:       password,
			METHOD:         method,
			PROTOCOL:       protocol,
			PROTOCOL_PARAM: protocolParam,
			OBFS:           obfs,
			OBFS_PARAM:     obfsParam,
			GROUP:          group,
			REMARK:         remark,
			CC:             cc,
			REGION:         region,
			VPN:            vpn,
		})
	}

	return result
}

func ToClash(Ssres []SsrStruct) string {
	var result = []string{"proxies:"}
	for _, Ssr := range Ssres {
		result = append(result, fmt.Sprintf("  - name: %s", Ssr.REMARK))
		result = append(result, fmt.Sprintf("    server: %s", Ssr.ADDRESS))
		result = append(result, fmt.Sprintf("    type: %s", Ssr.VPN))
		result = append(result, fmt.Sprintf("    port: %d", Ssr.PORT))
		result = append(result, fmt.Sprintf("    cipher: %s", Ssr.METHOD))
		result = append(result, fmt.Sprintf("    password: %s", Ssr.PASSWORD))
		result = append(result, fmt.Sprintf("    obfs: %s", Ssr.OBFS))
		result = append(result, fmt.Sprintf("    obfs-param: %s", Ssr.OBFS_PARAM))
		result = append(result, fmt.Sprintf("    protocol: %s", Ssr.PROTOCOL))
		result = append(result, fmt.Sprintf("    protocol-param: %s", Ssr.PROTOCOL_PARAM))
		result = append(result, "    udp: true")
	}

	return strings.Join(result[:], "\n")
}

func ToSingBox(ssrs []SsrStruct) string {
	var result []string

	for _, ssr := range ssrs {
		result = append(result, fmt.Sprintf(`
		{
			"type": "shadowsocksr",
			"tag": "%s",
			"server": "%s",
			"server_port": %d,
			"method": "%s",
			"password": "%s",
			"obfs": "%s",
			"obfs_param": "%s",
			"protocol": "%s",
			"protocol_param": "%s"
		}`, ssr.REMARK, ssr.ADDRESS, ssr.PORT, ssr.METHOD, ssr.PASSWORD, ssr.OBFS, ssr.OBFS_PARAM, ssr.PROTOCOL, ssr.PROTOCOL_PARAM))
	}

	return fmt.Sprintf(`
		{
			"outbounds": [%s]
		}`, strings.Join(result[:], ","))
}

func ToRaw(Ssres []SsrStruct) string {
	var result []string

	for _, Ssr := range Ssres {
		password := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(Ssr.PASSWORD)), "=")
		remarks := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(Ssr.REMARK)), "=")
		protoParam := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(Ssr.PROTOCOL_PARAM)), "=")
		obfsParam := strings.TrimRight(base64.StdEncoding.EncodeToString([]byte(Ssr.OBFS_PARAM)), "=")
		fmt.Println(Ssr.OBFS_PARAM)
		result = append(result, "ssr://"+base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%d:%s:%s:%s:%s/?remarks=%s&protoparam=%s&obfsparam=%s", Ssr.ADDRESS, Ssr.PORT, Ssr.PROTOCOL, Ssr.METHOD, Ssr.OBFS, password, remarks, protoParam, obfsParam))))
	}

	return strings.Join(result[:], "\n")
}
