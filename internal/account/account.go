package account

import (
	"github.com/LalatinaHub/LatinaSub-go/db"
)

func Get(filter string) []db.DBScheme {
	return db.New().Get(filter)
}
