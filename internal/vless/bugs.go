package vless

import (
	"math/rand"
)

func FillBugs(vlesses []VlessStruct, cdn, sni []string) []VlessStruct {
	for i := range vlesses {
		randSni := rand.Intn(len(sni))
		randCdn := rand.Intn(len(cdn))

		if vlesses[i].IS_CDN {
			vlesses[i].ADDRESS = cdn[randCdn]
		} else {
			vlesses[i].HOST = sni[randSni]
			vlesses[i].SNI = sni[randSni]
		}
	}

	return vlesses
}
