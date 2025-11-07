package storage

import (
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

	_, err = statement.Exec(0, s.ID)
	if err != nil {
		return err
	}
	statement.Close()

	// insert into outages table
	sql = "insert into outages values (null, ?, ?, '0000-00-00 00:00:00');"
	statement, err = c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(s.ID, time.Now().UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return statement.Close()
}

func endOutage(c *Connection, s *models.Site) error {
	// update site table
	sql := "UPDATE sites SET is_up = ? WHERE id = ?"
	statement, err := c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = statement.Exec(1, s.ID)
	if err != nil {
		return err
	}
	statement.Close()

	sql = "update outages set outage_end = ? where website_id = ? and outage_end = '0000-00-00 00:00:00'"
	statement, err = c.DB.Prepare(sql)
	if err != nil {
		return err
	}

	statement.Exec(time.Now().UTC().Format("2006-01-02 15:04:05"), s.ID)
	return statement.Close()

	// s.Logger.Info(s.URL + " is back up")

	// if secondsDown >= 35 {
	// 	go func() {
	// 		err := checkSMTPEnv()
	// 		if err != nil {
	// 			s.Logger.Error(err.Error())
	// 			return
	// 		}

	// 		s.Logger.Info("Sending website up email for " + s.URL)
	// 		m := gomail.NewMessage()
	// 		m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	// 		m.SetHeader("To", os.Getenv("EMAIL_TO"))
	// 		m.SetHeader("Subject", s.URL+" is back up")
	// 		m.SetBody("text/html", "<h1>"+s.URL+" is back online!</h1><p>It was down for "+strconv.Itoa(secondsDown)+" seconds.</p>")

	// 		port := os.Getenv("SMTP_PORT")
	// 		portInt, _ := strconv.Atoi(port)
	// 		d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
	// 			portInt,
	// 			os.Getenv("SMTP_USER"),
	// 			os.Getenv("SMTP_PASSWORD"))
	// 		if err := d.DialAndSend(m); err != nil {
	// 			s.Logger.Error(err)
	// 		}
	// 	}()
	// }
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
