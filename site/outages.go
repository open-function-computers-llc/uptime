package site

import (
	"fmt"
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
)

//Outage - a specific instance of a site going down
type Outage struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration float64   `json:"duration"`
}

func (s *Website) Outages(dbConn *storage.Connection) []Outage {
	outages := []Outage{}

	row, err := dbConn.DB.Query("SELECT outage_start, outage_end FROM outages WHERE website_id = ?", s.ID)
	if err != nil {
		fmt.Println(err)
		return outages
	}
	defer row.Close()
	for row.Next() {
		var outage Outage
		var start string
		var end string
		layout := "2006-01-02T15:04:05Z"
		row.Scan(&start, &end)

		outage.Start, _ = time.Parse(layout, start)
		outage.End, _ = time.Parse(layout, end)
		if !outage.End.IsZero() {
			outage.Duration = outage.End.Sub(outage.Start).Seconds()
		}

		outages = append(outages, outage)
	}
	return outages
}
