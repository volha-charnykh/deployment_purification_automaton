package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/xo/dburl"
)

func InitDb(dbUrl string) *sql.DB {
	db, err := dburl.Open(dbUrl)
	if err != nil {
		log.Fatalf("Error while connecting to database: %v", err)
	}
	return db
}
