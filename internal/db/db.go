package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DbPath string = "resources/db.sqlite"
	DbUrl  string = "https://github.com/LalatinaHub/LatinaSub/raw/main/result/db.sqlite"
)

type DB struct {
	connPool []sql.DB
}

func (x *DB) Init() {
	for i := 0; i < 10; i++ {
		db, err := sql.Open("sqlite3", DbPath)
		if err != nil {
			log.Fatal(err)
		}

		x.connPool = append(x.connPool, *db)
	}
}

func (x *DB) Connect() sql.DB {
	if len(x.connPool) < 1 {
		log.Fatal("Too much connection to database!")
	}

	conn := x.connPool[0]
	x.connPool = x.connPool[1:]
	return conn
}

func (x *DB) Close(conn sql.DB) {
	x.connPool = append(x.connPool, conn)
}

func (x DB) ConnLeft() int {
	return len(x.connPool)
}

var Database DB = DB{}
