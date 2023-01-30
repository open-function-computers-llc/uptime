package site

import (
	"os"
	"strconv"
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
	"gopkg.in/gomail.v2"
)

func (s *Website) beginOutage(dbConn *storage.Connection) {
	siteID := s.GetSiteID(dbConn)
	sql := "insert into outages values (null, ?, ?, '0000-00-00 00:00:00');"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		s.Logger.Error(err)
	}
	_, err = statement.Exec(siteID, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		s.Logger.Error(err)
	}
	statement.Close()
}

func (s *Website) endOutage(dbConn *storage.Connection, secondsDown int) {
	siteID := s.GetSiteID(dbConn)
	sql := "update outages set outage_end = ? where website_id = ? and outage_end = '0000-00-00 00:00:00'"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		s.Logger.Error(err)
	}
	statement.Exec(time.Now().Format("2006-01-02 15:04:05"), siteID)
	statement.Close()

	s.Logger.Info(s.URL + " is back up")

	if secondsDown >= 35 {
		go func() {
			err := checkSMTPEnv()
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			s.Logger.Info("Sending website up email for " + s.URL)
			m := gomail.NewMessage()
			m.SetHeader("From", os.Getenv("EMAIL_FROM"))
			m.SetHeader("To", os.Getenv("EMAIL_TO"))
			m.SetHeader("Subject", s.URL+" is back up")
			m.SetBody("text/html", "<h1>"+s.URL+" is back online!</h1><p>It was down for "+strconv.Itoa(secondsDown)+" seconds.</p>")

			port := os.Getenv("SMTP_PORT")
			portInt, _ := strconv.Atoi(port)
			d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
				portInt,
				os.Getenv("SMTP_USER"),
				os.Getenv("SMTP_PASSWORD"))
			if err := d.DialAndSend(m); err != nil {
				s.Logger.Error(err)
			}
		}()
	}
}
