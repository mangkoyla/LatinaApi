package vmess

type VmessStruct struct {
	ADDRESS          string `json:"address"`
	ALTER_ID         int    `json:"alterId"`
	PORT             int    `json:"port"`
	PASSWORD         string `json:"password"`
	SECURITY         string `json:"security"`
	HOST             string `json:"host"`
	TLS              bool   `json:"tls"`
	NETWORK          string `json:"network"`
	PATH             string `json:"path"`
	SKIP_CERT_VERIFY bool   `json:"skipCertVerify"`
	SNI              string `json:"sni"`
	REMARK           string `json:"remark"`
	IS_CDN           bool   `json:"cdn"`
	CC               string `json:"cc"`
	REGION           string `json:"region"`
	VPN              string `json:"vpn"`
}
