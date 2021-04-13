package site

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

// Website - a site that we will be checking
type Website struct {
	ID                   int
	URL                  string
	IsUp                 bool
	DB                   *storage.Connection
	Logger               *logrus.Logger
	standardWarningSent  bool
	emergencyWarningSent bool
}

// Create - Make a new instance of a Website struct
func Create(address string, dbConn *storage.Connection, logger *logrus.Logger) Website {
	w := Website{
		URL:    address,
		IsUp:   true,
		DB:     dbConn,
		Logger: logger,
	}
	logger.Info("Created Website:", address)

	siteDatabaseID := w.GetSiteID(dbConn)
	if siteDatabaseID == 0 {
		err := storage.AddSite(w.URL, dbConn)
		if err != nil {
			logger.Info("Couldn't add new site to DB:", err)
		}
	}
	return w
}

// Monitor - periodically make an HTTP GET request to the site's URL, and optionally log it in the database
func (s *Website) Monitor(shutdownChan *chan string) {
	go func() {
		secondsDown := 0

		for {
			select {
			case msg := <-*shutdownChan:
				if msg == s.URL {
					s.Logger.Info("Shutting down monitor for " + s.URL)
					return
				}

				s.Logger.Info("Site: " + s.URL + "passing url back to channel " + msg)
				*shutdownChan <- msg
			default:
				// nothing to do as the default, this is just here so that the
				// channel checking is non-blocking
			}

			statusCodeOrTimeoutValue := s.getStatusCode()
			// s.Logger.Info(s.URL+":", statusCode)
			if statusCodeOrTimeoutValue <= 99 {
				// fake status code returned because of slow response
				secondsDown += statusCodeOrTimeoutValue
				s.setSiteDown(s.DB, secondsDown)

				secondsDown += 5
				time.Sleep(time.Second * 5) // wait 5 seconds and try again
				continue
			}
			if statusCodeOrTimeoutValue == 200 {
				s.setSiteUp(s.DB, secondsDown)

				secondsDown = 0

				time.Sleep(time.Second * 30) // wait 30 seconds and try again
				continue
			}

			s.setSiteDown(s.DB, secondsDown)

			// wait 5 seconds and try again
			time.Sleep(time.Second * 5)
			secondsDown += 5
		}
	}()
}

func (s *Website) getStatusCode() int {
	if s.URL == "" {
		return 404
	}
	timeoutSeconds := 10

	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
			// s.Logger.Error(s.URL + " took too long to respond, and the URL was different!")
			return timeoutSeconds
		}
		s.Logger.Error(err.Error())
		return 500
	}
	req.Header.Set("User-Agent", "OFC_Uptime_Bot-version:"+os.Getenv("VERSION"))

	resp, err := client.Do(req)
	if err != nil {
		if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
			// s.Logger.Error(s.URL + " took too long to respond!")
			return timeoutSeconds
		}
		s.Logger.Error(err.Error())
		return 500
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func (s *Website) setSiteUp(dbConn *storage.Connection, secondsDown int) {
	if !s.IsUp {
		s.endOutage(dbConn, secondsDown)
	}
	s.IsUp = true
	s.emergencyWarningSent = false
	s.standardWarningSent = false

	sql := "UPDATE sites SET last_checked = ?, is_up = ? WHERE url = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		s.Logger.Error(err)
	}
	_, err = statement.Exec(time.Now().Format("2006-01-02 15:04:05"), 1, s.URL)
	if err != nil {
		s.Logger.Error(err)
	}
	statement.Close()
}

func (s *Website) setSiteDown(dbConn *storage.Connection, secondsDown int) {
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

	s.Logger.Info(s.URL + " has been down for at least " + strconv.Itoa(secondsDown) + " second(s)")

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

	if secondsDown >= 30 {
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

// FindWebsiteByID - Find a site in the DB by it's ID
func FindWebsiteByID(id int, dbConn *storage.Connection, logger *logrus.Logger) (Website, error) {
	s := Website{}
	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE id = ?", id)
	if err != nil {
		logger.Error(err)
		return s, err
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
		return site, nil
	}
	return s, nil
}

// FindWebsiteByURL - Find a site in the DB by it's URL
func FindWebsiteByURL(url string, dbConn *storage.Connection, logger *logrus.Logger) (Website, error) {
	s := Website{}
	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE url = ?", url)
	if err != nil {
		logger.Error(err)
		return s, err
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
		return site, nil
	}
	return s, errors.New("Site was not found when looping over DB records")
}

// Destroy - delete a site from the DB and close down the monitoring routine
func (s *Website) Destroy(c *chan string, dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	statement.Exec(1, s.ID)
	statement.Close()
	go func() {
		*c <- s.URL
	}()
}

// Purge - delete a site from the DB and close down the monitoring routine
func (s *Website) Purge(dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "DELETE FROM sites WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	statement.Exec(s.ID)
	statement.Close()
}

// Destroy - delete a site from the DB and close down the monitoring routine
func (s *Website) Restore(dbConn *storage.Connection, logger *logrus.Logger) {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		logger.Error(err)
	}
	s.Logger = logger
	statement.Exec(0, s.ID)
	logger.Info("Restored Website:", s.URL)
	statement.Close()
}

func (s *Website) getDomain() string {
	domain := strings.Replace(s.URL, "http://www.", "", 1)
	domain = strings.Replace(domain, "https://www.", "", 1)
	domain = strings.Replace(domain, "http://", "", 1)
	domain = strings.Replace(domain, "https://", "", 1)

	if domain[len(domain)-1:] == "/" {
		domain = domain[:len(domain)-1]
	}
	return domain
}
