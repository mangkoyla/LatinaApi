package vmess

import (
	"math/rand"
)

func FillBugs(vmesses []VmessStruct, cdn, sni []string) []VmessStruct {
	for i := range vmesses {
		randSni := rand.Intn(len(sni))
		randCdn := rand.Intn(len(cdn))

		if vmesses[i].IS_CDN {
			vmesses[i].ADDRESS = cdn[randCdn]
		} else {
			vmesses[i].HOST = sni[randSni]
			vmesses[i].SNI = sni[randSni]
		}
	}

	return vmesses
}
