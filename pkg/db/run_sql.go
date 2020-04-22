package db

import (
	"database/sql"
	"fmt"
	"os"

	//
	_ "github.com/lib/pq"
)

// RunSQL run sql statement
func RunSQL(sqlStatement string) (*sql.Rows, error) {
	password := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf(
		"host=db user=wazuplay password=%s dbname=wazuplay sslmode=disable",
		password,
	)
	db, err := sql.Open("postgres", connectionString)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(sqlStatement)
	return rows, err
}
