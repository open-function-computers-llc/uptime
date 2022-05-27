package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// FindWebsiteByID - Find a site in the DB by it's ID
func FindWebsiteByID(id int, dbConn *storage.Connection, logger *logrus.Logger) (Website, error) {
	s := Website{}
	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE id = ?", id)
	if err != nil {
		logger.Error(err)
		return s, err
	}
	defer row.Close()

	for row.Next() {
		var id int
		var url string
		var isUp int
		row.Scan(&id, &url, &isUp)
		site := Website{
			URL:    url,
			ID:     id,
			IsUp:   isUp == 1,
			DB:     dbConn,
			Logger: logger,
		}
		return site, nil
	}
	return s, nil
}
