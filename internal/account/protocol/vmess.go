package protocol

type VmessStruct struct {
	Add            string `json:"add"`
	Aid            int    `json:"aid"`
	Host           string `json:"host"`
	Id             string `json:"id"`
	Net            string `json:"net"`
	Path           string `json:"path"`
	Port           uint16 `json:"port"`
	Ps             string `json:"ps"`
	Tls            string `json:"tls"`
	Security       string `json:"security"`
	SkipCretVerify bool   `json:"skip-cert-verify"`
	Sni            string `json:"sni"`
}
