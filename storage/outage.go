package storage

import (
	"errors"
	"strconv"
	"time"

	"github.com/open-function-computers-llc/uptime/models"
)

func beginOutage(c *Connection, s *models.Site) error {
	// update site table
	sql := "UPDATE sites SET is_up = ? WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(0, s.ID)
	if err != nil {
		return err
	}

	// insert into outages table
	sql = "insert into outages values (null, ?, ?, '0000-00-00 00:00:00');"
	statement2, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement2.Close()

	result, err := statement2.Exec(s.ID, time.Now().UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	outageID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.Logger.Info("Created outage for site " + s.URL)
	c.Logger.Info("Outage ID: " + strconv.Itoa(int(outageID)))

	return nil
}

func endOutage(c *Connection, s *models.Site) error {
	c.Logger.Info("Bringing back up site with ID: " + strconv.Itoa(s.ID))

	// update site table
	sql := "UPDATE sites SET is_up = ? WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(1, s.ID)
	if err != nil {
		return err
	}

	// find oldest outage
	rows, err := c.DB.Query("SELECT id FROM outages WHERE website_id = ? and outage_end = '0000-00-00 00:00:00' ORDER BY outage_start ASC LIMIT 1", s.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	outageID := 0

	for rows.Next() {
		err := rows.Scan(&outageID)
		if err != nil {
			return err
		}
	}

	if outageID == 0 {
		return errors.New("couldn't find an outage to finish... wtf?")
	}

	sql = "UPDATE outages SET outage_end = ? WHERE id = ?"
	statement2, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}
	defer statement2.Close()

	_, err = statement2.Exec(time.Now().UTC().Format("2006-01-02 15:04:05"), outageID)
	return err
}

func SiteOutages(siteID int, c *Connection) ([]*models.Outage, error) {
	rows, err := c.DB.Query(`
        SELECT id, website_id, outage_start, outage_end
        FROM outages
        WHERE website_id = ?
		ORDER BY outage_start
    `, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	output := []*models.Outage{}

	for rows.Next() {
		outage := &models.Outage{}

		err := rows.Scan(&outage.ID, &outage.WebsiteID, &outage.Start, &outage.End)
		if err != nil {
			return nil, err
		}

		if outage.End.IsZero() {
			outage.End = time.Now().UTC()
		}

		outage.Duration = int(outage.End.Sub(outage.Start).Seconds())

		output = append(output, outage)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return output, nil
}

func SiteMostRecentOutageDuration(siteID int, c *Connection) int {
	outages, _ := SiteOutages(siteID, c)
	if len(outages) == 0 {
		return 0
	}
	mostRecentOutage := outages[len(outages)-1]

	return mostRecentOutage.Duration
}
