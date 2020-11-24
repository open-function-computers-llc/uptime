package storage

import "database/sql"

type Connection struct {
	DB *sql.DB
}
