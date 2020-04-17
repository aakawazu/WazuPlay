package db

import (
	"database/sql"
	//
	_ "github.com/lib/pq"
)

// ConnectToDB connect to db
func ConnectToDB() *sql.DB {
	db, _ := sql.Open("postgres", "host=db user=wazuplay password=test")
	return db
}
