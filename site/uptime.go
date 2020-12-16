package site

import (
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
)

// CalcUptime - Calculate the percentage of uptime over the last X days
func (s *Website) CalcUptime(days int, dbConn *storage.Connection) float64 {
	outages := s.Outages(dbConn)

	if len(outages) == 0 {
		return 1
	}
	// here are the total number of seconds possible in the requested duration
	totalSecondsInRange := float64(days * 60 * 60 * 24)
	uptime := totalSecondsInRange
	now := time.Now()
	rangeBegin := now.AddDate(0, 0, days*-1)

	for _, outage := range outages {
		if outage.End.After(rangeBegin) {
			uptime -= outage.Duration
		}
	}
	return uptime / totalSecondsInRange
}
