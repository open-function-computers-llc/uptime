package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func setUpStorage() (*sql.DB, error) {
	dbConn, err := sql.Open("mysql", "lapubell:genius@tcp(localhost:3306)/ofc_uptime?parseTime=true")
	return dbConn, err
}
