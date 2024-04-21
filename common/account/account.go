package account

import (
	"github.com/mangkoyla/LatinaSub-go/db"
)

func Get(filter string) []db.DBScheme {
	return db.New().Get(filter)
}
