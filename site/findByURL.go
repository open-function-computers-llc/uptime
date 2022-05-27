package site

import (
	"errors"

	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
)

// FindWebsiteByURL - Find a site in the DB by it's URL
func FindWebsiteByURL(url string, dbConn *storage.Connection, logger *logrus.Logger) (Website, error) {
	s := Website{}
	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE url like ?", "%"+url+"%")
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
			URL:  url,
			ID:   id,
			IsUp: isUp == 1,
		}
		return site, nil
	}
	return s, errors.New("Site was not found when looping over DB records")
}
