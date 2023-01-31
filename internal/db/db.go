package db

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DbName string = ""
	DbPath string = "resources/databases/"
	DbUrl  string = "https://github.com/LalatinaHub/LatinaSub/raw/main/result/db.sqlite"
)

func selectDb() string {
	var (
		latestDb int
	)

	files, _ := ioutil.ReadDir("resources/databases/")
	for _, file := range files {
		info, _ := os.Stat(DbPath + file.Name())
		modTime := info.ModTime().Format("20060102150405")
		newDb, _ := strconv.Atoi(modTime)

		if newDb > latestDb {
			latestDb = newDb

			DbName = info.Name()
		}
	}

	return DbPath + DbName
}

func Connect() sql.DB {
	db, err := sql.Open("sqlite3", selectDb())
	if err != nil {
		log.Fatal(err)
	}

	return *db
}
