package storage

import (
	"os"
	"strconv"
	"time"

	"github.com/open-function-computers-llc/uptime/email"
	"github.com/open-function-computers-llc/uptime/models"
	"github.com/sirupsen/logrus"
)

func CreateSite(url, meta string, c *Connection) (*models.Site, error) {
	sql := "INSERT INTO sites (url, meta, is_up, is_deleted, created_at, last_checked) VALUES (?, ?, 1, 0, ?, ?)"

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return nil, err
	}

	_, err = statement.Exec(url, meta, time.Now().UTC().Format("2006-01-02 15:04:05"), time.Now().UTC().Format("2006-01-02 15:04:05"))

	return nil, err
}

func UpdateSite(id int, url, meta string, dbConn *Connection) error {
	sql := "UPDATE sites SET url = ?, meta = ? WHERE id = ?"

	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(url, meta, id)

	return err
}

// GetSites - get all the sites out of the storage connection
func GetSites(c *Connection) ([]*models.Site, error) {
	sites := []*models.Site{}

	query := "SELECT s.id, s.url, s.meta, s.is_up, s.is_deleted, "

	// calculate uptimes
	query += `
	-- 1 day
	ROUND(
      100 * (
        1 - IFNULL(
          SUM(
            CASE
              WHEN
                (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END)
                  >= (UTC_TIMESTAMP() - INTERVAL 1 DAY)
                AND o.outage_start <= UTC_TIMESTAMP()
              THEN
                TIMESTAMPDIFF(SECOND,
                  GREATEST(o.outage_start, UTC_TIMESTAMP() - INTERVAL 1 DAY),
                  LEAST(
                    /* effective end again */
                    (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END),
                    UTC_TIMESTAMP()
                  )
                )
              ELSE 0
            END
          ) / 86400
        , 0)
      )
    , 2) AS uptime_24h,

    -- 7 days
    ROUND(
      100 * (
        1 - IFNULL(
          SUM(
            CASE
              WHEN
                (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END)
                  >= (UTC_TIMESTAMP() - INTERVAL 7 DAY)
                AND o.outage_start <= UTC_TIMESTAMP()
              THEN
                TIMESTAMPDIFF(SECOND,
                  GREATEST(o.outage_start, UTC_TIMESTAMP() - INTERVAL 7 DAY),
                  LEAST(
                    (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END),
                    UTC_TIMESTAMP()
                  )
                )
              ELSE 0
            END
          ) / 604800
        , 0)
      )
    , 2) AS uptime_7d,

    -- 30 days
    ROUND(
      100 * (
        1 - IFNULL(
          SUM(
            CASE
              WHEN
                (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END)
                  >= (UTC_TIMESTAMP() - INTERVAL 30 DAY)
                AND o.outage_start <= UTC_TIMESTAMP()
              THEN
                TIMESTAMPDIFF(SECOND,
                  GREATEST(o.outage_start, UTC_TIMESTAMP() - INTERVAL 30 DAY),
                  LEAST(
                    (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END),
                    UTC_TIMESTAMP()
                  )
                )
              ELSE 0
            END
          ) / 2592000
        , 0)
      )
    , 2) AS uptime_30d,

    -- 60 days
    ROUND(
      100 * (
        1 - IFNULL(
          SUM(
            CASE
              WHEN
                (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END)
                  >= (UTC_TIMESTAMP() - INTERVAL 60 DAY)
                AND o.outage_start <= UTC_TIMESTAMP()
              THEN
                TIMESTAMPDIFF(SECOND,
                  GREATEST(o.outage_start, UTC_TIMESTAMP() - INTERVAL 60 DAY),
                  LEAST(
                    (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END),
                    UTC_TIMESTAMP()
                  )
                )
              ELSE 0
            END
          ) / 5184000
        , 0)
      )
    , 2) AS uptime_60d,

    -- 90 days
    ROUND(
      100 * (
        1 - IFNULL(
          SUM(
            CASE
              WHEN
                (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END)
                  >= (UTC_TIMESTAMP() - INTERVAL 90 DAY)
                AND o.outage_start <= UTC_TIMESTAMP()
              THEN
                TIMESTAMPDIFF(SECOND,
                  GREATEST(o.outage_start, UTC_TIMESTAMP() - INTERVAL 90 DAY),
                  LEAST(
                    (CASE WHEN o.outage_end = '0000-00-00 00:00:00' OR o.outage_end IS NULL THEN UTC_TIMESTAMP() ELSE o.outage_end END),
                    UTC_TIMESTAMP()
                  )
                )
              ELSE 0
            END
          ) / 7776000
        , 0)
      )
    , 2) AS uptime_90d
	`

	query += " FROM sites s LEFT JOIN outages o ON o.website_id = s.id GROUP BY s.id "

	rows, err := c.DB.Query(query)
	if err != nil {
		return sites, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var site models.Site
		var isUp int
		var isDeleted int

		err = rows.Scan(
			&site.ID,
			&site.URL,
			&site.Meta,
			&isUp,
			&isDeleted,
			&site.Uptime_1day,
			&site.Uptime_7day,
			&site.Uptime_30day,
			&site.Uptime_60day,
			&site.Uptime_90day,
		)

		if err != nil {
			return sites, err
		}

		site.IsUp = isUp == 1
		site.IsDeleted = isDeleted == 1
		sites = append(sites, &site)
		count++
	}

	return sites, nil
}

func FindSiteByID(c *Connection, id int) (models.Site, error) {
	s := models.Site{}

	row, err := c.DB.Query("SELECT id, url, meta, is_up FROM sites WHERE id = ?", id)
	if err != nil {
		return s, err
	}
	defer row.Close()

	for row.Next() {
		var id int
		var url string
		var isUp int
		var meta string
		row.Scan(&id, &url, &meta, &isUp)
		site := models.Site{
			URL:  url,
			ID:   id,
			IsUp: isUp == 1,
		}
		return site, nil
	}
	return s, nil
}

// Destroy - delete a site from the DB and close down the monitoring routine
func DeleteSite(c *Connection, id int) error {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(1, id)
	if err != nil {
		return err
	}

	return statement.Close()
}

// Destroy - delete a site from the DB and close down the monitoring routine
func RestoreSite(c *Connection, id int) error {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(0, id)
	if err != nil {
		return err
	}

	return statement.Close()
}

func Monitor(s *models.Site, c *Connection) {
	if s.IsMonitoring {
		return
	}

	s.IsMonitoring = true
	okSeconds, _ := strconv.Atoi(os.Getenv("INTERVAL_OK_RECHECK"))
	errorSeconds, _ := strconv.Atoi(os.Getenv("INTERVAL_ERROR_RECHECK"))
	dangerSeconds, _ := strconv.Atoi(os.Getenv("DANGER_SECONDS"))
	emergencySeconds, _ := strconv.Atoi(os.Getenv("EMERGENCY_SECONDS"))
	c.Logger.Info("starting monitor for URL:", s.URL)

	secondsDown := 0

	for {
		if !s.IsMonitoring {
			c.Logger.Info("shutting down monitor for URL:", s.URL)
			break
		}

		c.Logger.Info("checking url: ", s.URL)
		UpdateLastChecked(s.ID, c)

		resCode, err := s.GetStatusCodeAndError()
		if resCode == 200 && err == nil {
			go SetSiteUp(c, s)
			time.Sleep(time.Second * time.Duration(okSeconds))
			continue
		}

		go SetSiteDown(c, s, c.Logger)

		if secondsDown >= dangerSeconds {
			webhooks, _ := SiteWebhooks(s.ID, c)
			go processStandardActions(secondsDown, s, webhooks)
		}
		if secondsDown >= emergencySeconds {
			webhooks, _ := SiteWebhooks(s.ID, c)
			go processEmergencyActions(secondsDown, s, webhooks)
		}

		// wait X seconds and try again
		secondsDown += errorSeconds
		time.Sleep(time.Second * time.Duration(errorSeconds))
	}
}

func StopMonitoringSite(s *models.Site, logger *logrus.Logger) {
	s.IsMonitoring = false
}

func SetSiteUp(c *Connection, s *models.Site) error {
	s.EmergencyWarningSent = false
	s.StandardWarningSent = false

	if !s.IsUp {
		s.IsUp = true
		c.Logger.Info("bringing this site back up: ", s)
		go email.Send(s.URL+" is back up", buildHTMLMessage(s.URL, s.Meta, 0), false)
		return endOutage(c, s)
	}
	return nil
}

func SetSiteDown(c *Connection, s *models.Site, logger *logrus.Logger) error {
	if s.IsUp {
		s.IsUp = false
		logger.Info("site is now down! ", s)
		return beginOutage(c, s)
	}

	return nil
}

func UpdateLastChecked(id int, c *Connection) error {

	sql := "UPDATE sites SET last_checked = ? WHERE id = ?"

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(time.Now().UTC().Format("2006-01-02 15:04:05"), id)

	return err
}

func PurgeSite(id int, c *Connection) error {
	// purge outages
	sql := "DELETE FROM outages WHERE website_id = ?"

	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	// purge webhooks
	sql = "DELETE FROM webhooks WHERE website_id = ?"

	statement, err = c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	// purge actual site
	sql = "DELETE FROM sites WHERE id = ?"

	statement, err = c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(id)

	return err
}
