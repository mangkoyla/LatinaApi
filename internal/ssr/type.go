package ssr

type SsrStruct struct {
	ADDRESS        string `json:"address"`
	PORT           int    `json:"port"`
	PASSWORD       string `json:"password"`
	METHOD         string `json:"method"`
	PROTOCOL       string `json:"protocol"`
	PROTOCOL_PARAM string `json:"protocolParam"`
	OBFS           string `json:"obfs"`
	OBFS_PARAM     string `json:"obfsParam"`
	GROUP          string `json:"group"`
	REMARK         string `json:"remark"`
	CC             string `json:"cc"`
	REGION         string `json:"region"`
	VPN            string `json:"vpn"`
}
