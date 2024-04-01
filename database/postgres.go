package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func StartDatabase() (*sql.DB, error) {
	var host = "postgres"
	var port = 5432
	var user = "postgres"
	var password = "admin"
	var dbname = "postgres"
	var sslmode = "disable"

	var dbinfo = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var db, err = sql.Open("postgres", dbinfo)
	log.Println("staging dbinfo", dbinfo)

	if err != nil {
		log.Fatalf("Failed to connect to postgres database: %v", err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(7)

	return db, nil
}
