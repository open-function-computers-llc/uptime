package storage

import (
	"fmt"
	"time"
)

func OutagesByDay(c *Connection) (map[string]int, error) {
	output := map[string]int{}

	sql := `SELECT DATE(outage_start) AS outage_day,
                   COUNT(*) AS total_outages
                   FROM outages
                   GROUP BY DATE(outage_start)
                   ORDER BY outage_day;`

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return output, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		return output, err
	}
	defer rows.Close()

	// Temporary map to hold raw dates from DB before normalization
	rawDates := make(map[string]int)

	for rows.Next() {
		var outageDay string
		var totalOutages int
		err := rows.Scan(&outageDay, &totalOutages)
		if err != nil {
			return output, err
		}
		rawDates[outageDay] = totalOutages
	}
	if err = rows.Err(); err != nil {
		return output, err
	}

	// 1. Normalize existing keys and find date range
	var startDate, endDate time.Time
	firstKey := true

	for key, count := range rawDates {
		var t time.Time
		// Try parsing the various formats DATE() might return
		if t, err = time.Parse(time.RFC3339, key); err != nil {
			if t, err = time.Parse("2006-01-02", key); err != nil {
				return output, fmt.Errorf("failed to parse date key: %s", key)
			}
		}

		// Truncate to midnight
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

		// Store with consistent key "YYYY-MM-DD"
		normalizedKey := t.Format("2006-01-02")
		output[normalizedKey] = count

		if firstKey {
			startDate = t
			endDate = t
			firstKey = false
		} else {
			if t.Before(startDate) {
				startDate = t
			}
			if t.After(endDate) {
				endDate = t
			}
		}
	}

	// 2. Fill in missing days with 0
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateKey := currentDate.Format("2006-01-02")
		if _, exists := output[dateKey]; !exists {
			output[dateKey] = 0
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return output, nil
}

func OutageDurationsByDay(c *Connection) (map[string]int, error) {
	output := map[string]int{}

	sql := `SELECT DATE(outage_start) AS outage_day,
			       SUM(TIMESTAMPDIFF(SECOND, outage_start, outage_end)) AS total_duration_seconds
				   FROM outages
				   GROUP BY DATE(outage_start)
				   ORDER BY outage_day;`

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return output, err
	}
	defer statement.Close()

	rows, err := statement.Query()
	if err != nil {
		return output, err
	}
	defer rows.Close()

	// Temporary map to hold raw dates from DB before normalization
	rawDates := make(map[string]int)

	for rows.Next() {
		var outageDay string
		var outageDurationSeconds int
		err := rows.Scan(&outageDay, &outageDurationSeconds)
		if err != nil {
			return output, err
		}
		rawDates[outageDay] = outageDurationSeconds
	}
	if err = rows.Err(); err != nil {
		return output, err
	}

	// 1. Normalize existing keys and find date range
	var startDate, endDate time.Time
	firstKey := true

	for key, count := range rawDates {
		var t time.Time
		// Try parsing the various formats DATE() might return
		if t, err = time.Parse(time.RFC3339, key); err != nil {
			if t, err = time.Parse("2006-01-02", key); err != nil {
				return output, fmt.Errorf("failed to parse date key: %s", key)
			}
		}

		// Truncate to midnight
		t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

		// Store with consistent key "YYYY-MM-DD"
		normalizedKey := t.Format("2006-01-02")
		output[normalizedKey] = count

		if firstKey {
			startDate = t
			endDate = t
			firstKey = false
		} else {
			if t.Before(startDate) {
				startDate = t
			}
			if t.After(endDate) {
				endDate = t
			}
		}
	}

	// 2. Fill in missing days with 0
	currentDate := startDate
	for !currentDate.After(endDate) {
		dateKey := currentDate.Format("2006-01-02")
		if _, exists := output[dateKey]; !exists {
			output[dateKey] = 0
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return output, nil
}
