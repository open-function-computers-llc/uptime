package site

import (
	"os"
	"strconv"
	"time"

	"github.com/open-function-computers-llc/uptime/notifications"
	"github.com/open-function-computers-llc/uptime/storage"
	"gopkg.in/gomail.v2"
)

func (s *Website) setSiteDown(dbConn *storage.Connection, secondsDown int, statusCode int) {
	if s.IsUp {
		s.beginOutage(dbConn)
	}
	s.IsUp = false
	sql := "UPDATE sites SET last_checked = ?, is_up = ? WHERE url = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		s.Logger.Error(err)
	}
	_, err = statement.Exec(time.Now().Format("2006-01-02 15:04:05"), 0, s.URL)
	if err != nil {
		s.Logger.Error(err)
	}
	statement.Close()

	s.Logger.Info(s.URL + " has been down for at least " + strconv.Itoa(secondsDown) + " second(s), current code: " + strconv.Itoa(statusCode))

	if secondsDown >= 60 && !s.standardWarningSent {
		go func() {
			err := checkSMTPEnv()
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			s.Logger.Info("Sending standard down email for " + s.URL)
			m := gomail.NewMessage()
			m.SetHeader("From", os.Getenv("EMAIL_FROM"))
			m.SetHeader("To", os.Getenv("EMAIL_TO"))
			m.SetHeader("Subject", s.URL+" is down!")
			m.SetBody("text/html", "<h1>"+s.URL+" is down!</h1><p>It has been down for at least "+strconv.Itoa(secondsDown)+" seconds. Better go check things out...</p>")

			port := os.Getenv("SMTP_PORT")
			portInt, _ := strconv.Atoi(port)
			d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
				portInt,
				os.Getenv("SMTP_USER"),
				os.Getenv("SMTP_PASSWORD"))
			if err := d.DialAndSend(m); err != nil {
				s.Logger.Error(err)
			}
			s.standardWarningSent = true
		}()
	}

	if secondsDown >= 180 && !s.emergencyWarningSent {
		go notifications.SendHTTPRequest(s.URL)
		go func() {
			err := checkSMTPEnv()
			if err != nil {
				s.Logger.Error(err.Error())
				return
			}

			s.Logger.Info("Sending emergency down email for " + s.URL)
			m := gomail.NewMessage()
			m.SetHeader("From", os.Getenv("EMAIL_FROM"))
			m.SetHeader("To", os.Getenv("EMAIL_TO"))
			m.SetHeader("Subject", "OFCO.911 - "+s.URL+" is down!")
			m.SetBody("text/html", "<h1>"+s.URL+" is down!</h1><p>It has been down for at least "+strconv.Itoa(secondsDown)+" seconds. Better go check things out...</p>")

			port := os.Getenv("SMTP_PORT")
			portInt, _ := strconv.Atoi(port)
			d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
				portInt,
				os.Getenv("SMTP_USER"),
				os.Getenv("SMTP_PASSWORD"))
			if err := d.DialAndSend(m); err != nil {
				s.Logger.Error(err)
			}
			s.emergencyWarningSent = true
		}()
	}
}
