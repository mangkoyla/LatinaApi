package account

import (
	"math/rand"

	"github.com/LalatinaHub/LatinaSub-go/db"
)

func PopulateBugs(accounts []db.DBScheme, cdn, sni []string) []db.DBScheme {
	for i, account := range accounts {
		randSni := rand.Intn(len(sni))
		randCdn := rand.Intn(len(cdn))

		switch account.ConnMode {
		case "cdn":
			accounts[i].Server = cdn[randCdn]
		case "sni":
			accounts[i].SNI = sni[randSni]
			accounts[i].Host = sni[randSni]
		}
	}

	return accounts
}
