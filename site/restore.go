package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Restore - restore a site from the DB and close down the monitoring routine
func (s *Website) Restore(dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	s.Logger = logger
	statement.Exec(0, s.ID)
	logger.Info("Restored Website:", s.URL)
	statement.Close()
}
