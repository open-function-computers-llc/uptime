package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// Destroy - delete a site from the DB and close down the monitoring routine
func (s *Website) Destroy(c *chan string, dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	statement.Exec(1, s.ID)
	statement.Close()
	go func() {
		*c <- s.URL
	}()
}
