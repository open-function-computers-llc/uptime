package main

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func setUpStorage() (*sql.DB, error) {
	dsn, err := createDSN()
	if err != nil {
		return nil, err
	}
	dbConn, err := sql.Open("mysql", dsn)
	return dbConn, err
}

func createDSN() (string, error) {
	// required fields
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")

	// this is pretty much always 3306, so this one can be optional
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	if dbName == "" || dbHost == "" || dbUser == "" || dbPass == "" {
		return "", errors.New("missing required ENV. Please check your db creds and try again")
	}

	return dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true", nil
}
