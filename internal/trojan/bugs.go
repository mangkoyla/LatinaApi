package trojan

import (
	"math/rand"
)

func FillBugs(trojans []TrojanStruct, cdn, sni []string) []TrojanStruct {
	for i := range trojans {
		randSni := rand.Intn(len(sni))
		randCdn := rand.Intn(len(cdn))

		if trojans[i].IS_CDN {
			trojans[i].ADDRESS = cdn[randCdn]
		} else {
			trojans[i].HOST = sni[randSni]
			trojans[i].SNI = sni[randSni]
		}
	}

	return trojans
}
