package ssr

import (
	"fmt"
	"math/rand"
	"regexp"
)

func FillBugs(ssrs []SsrStruct, sni []string) []SsrStruct {
	tlsObfs, _ := regexp.Compile("tls")
	for i := range ssrs {
		randSni := rand.Intn(len(sni))

		if tlsObfs.MatchString(ssrs[i].OBFS) {
			ssrs[i].OBFS_PARAM = fmt.Sprintf("obfs=tls;obfs-host=%s", sni[randSni])
		} else {
			ssrs[i].OBFS_PARAM = fmt.Sprintf("obfs=http;obfs-host=%s", sni[randSni])
		}
	}

	return ssrs
}
