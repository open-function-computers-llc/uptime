package site

import (
	"github.com/open-function-computers-llc/uptime/storage"
)

// GetSiteID - the the ID of the current website by it's URL
func (s *Website) GetSiteID(dbConn *storage.Connection) int {
	var siteID int
	row, err := dbConn.DB.Query("SELECT id FROM sites WHERE url = '" + s.URL + "'")
	if err != nil {
		s.Logger.Error(err)
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&siteID)
	}
	return siteID
}
