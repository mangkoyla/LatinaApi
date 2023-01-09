package trojan

type TrojanStruct struct {
	ADDRESS          string `json:"address"`
	PORT             int    `json:"port"`
	PASSWORD         string `json:"password"`
	SECURITY         string `json:"security"`
	HOST             string `json:"host"`
	TLS              bool   `json:"tls"`
	NETWORK          string `json:"network"`
	PATH             string `json:"path"`
	SERVICE_NAME     string `json:"serviceName"`
	MODE             string `json:"mode"`
	SKIP_CERT_VERIFY bool   `json:"skipCertVerify"`
	SNI              string `json:"sni"`
	FLOW             string `json:"flow"`
	LEVEL            int    `json:"level"`
	METHOD           string `json:"method"`
	OTA              bool   `json:"ota"`
	REMARK           string `json:"remark"`
	IS_CDN           bool   `json:"cdn"`
	CC               string `json:"cc"`
	REGION           string `json:"region"`
	VPN              string `json:"vpn"`
}
