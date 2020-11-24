package site

import (
	"fmt"

	"github.com/open-function-computers-llc/uptime/storage"
)

// GetSites - get all the sites out of the storage connection
func GetSites(dbConn *storage.Connection) map[int]Website {
	sites := map[int]Website{}

	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE is_deleted = 0")
	if err != nil {
		fmt.Println(err)
		return sites
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
		sites[id] = site
	}
	return sites
}
