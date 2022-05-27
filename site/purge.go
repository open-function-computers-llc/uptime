package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Purge - delete a site from the DB and close down the monitoring routine
func (s *Website) Purge(dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "DELETE FROM sites WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	statement.Exec(s.ID)
	statement.Close()
}
