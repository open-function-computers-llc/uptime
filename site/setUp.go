package site

import (
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *Website) setSiteUp(dbConn *storage.Connection, secondsDown int) {
	if !s.IsUp {
		s.endOutage(dbConn, secondsDown)
	}
	s.IsUp = true
	s.emergencyWarningSent = false
	s.standardWarningSent = false

	sql := "UPDATE sites SET last_checked = ?, is_up = ? WHERE url = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		s.Logger.Error(err)
	}
	_, err = statement.Exec(time.Now().Format("2006-01-02 15:04:05"), 1, s.URL)
	if err != nil {
		s.Logger.Error(err)
	}
	statement.Close()
}
