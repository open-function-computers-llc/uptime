package storage

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

// Connection - A struct containing all the stuff needed for a storage instance
type Connection struct {
	DB     *sql.DB
	Logger *logrus.Logger
}
