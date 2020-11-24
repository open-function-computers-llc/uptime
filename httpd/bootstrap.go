package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func setUpStorage() (*sql.DB, error) {
	sqliteDB, err := sql.Open("sqlite3", "./database")
	return sqliteDB, err
}

func shutDownStorage(db *sql.DB) {
	db.Close()
}
