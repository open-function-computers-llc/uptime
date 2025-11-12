package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// connection - A struct containing all the stuff needed for a storage instance
type Connection struct {
	DB     *sql.DB
	Logger *logrus.Logger
}

func NewConnection(logger *logrus.Logger) (*Connection, error) {
	c := Connection{
		Logger: logger,
	}

	db, err := setUpStorage()
	if err != nil {
		return nil, err
	}

	c.DB = db

	return &c, nil
}

func setUpStorage() (*sql.DB, error) {
	dsn, err := createDSN()
	if err != nil {
		return nil, err
	}

	dbConn, err := sql.Open(os.Getenv("DB_TYPE"), dsn)
	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, err
	}

	// Force all timestamps and NOW()/UTC_TIMESTAMP() to use UTC
	_, err = dbConn.Exec("SET time_zone = '+00:00'")
	if err != nil {
		return nil, fmt.Errorf("failed to set UTC timezone: %w", err)
	}

	return dbConn, nil
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

func (c *Connection) Close() error {
	return c.DB.Close()
}
