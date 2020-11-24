package site

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/open-function-computers-llc/uptime/storage"
)

// Website - a site that we will be checking
type Website struct {
	ID   int
	URL  string
	IsUp bool
	DB   *storage.Connection
}

// Create - Make a new instance of a Website struct
func Create(address string, dbConn *storage.Connection) Website {
	w := Website{
		URL:  address,
		IsUp: true,
		DB:   dbConn,
	}
	log.Println("Created Website:", address)

	siteDatabaseID := w.GetSiteID(dbConn)
	if siteDatabaseID == 0 {
		err := storage.AddSite(w.URL, dbConn)
		if err != nil {
			log.Println("Couldn't add new site to DB:", err)
		}
	}
	return w
}

// Monitor - periodically make an HTTP GET request to the site's URL, and optionally log it in the database
func (s *Website) Monitor(shutdownChan chan string) {
	go func() {
		for {
			select {
			case msg := <-shutdownChan:
				if msg == s.URL {
					log.Println("Shutting down monitor for " + s.URL)
					return
				} else {
					log.Println("Site: " + s.URL + "passing url back to channel " + msg)
					shutdownChan <- msg
				}
			default:
				// nothing to do as the default, this is just here so that the
				// channel checking is non-blocking
			}

			statusCode := s.getStatusCode()
			log.Println(s.URL+":", statusCode)
			if statusCode == 200 {
				s.setSiteUp(s.DB)
				time.Sleep(time.Second * 15)
				continue
			}

			s.setSiteDown(s.DB)
			time.Sleep(time.Second * 1)
		}
	}()
}

func (s *Website) getStatusCode() int {
	if s.URL == "" {
		return 404
	}

	resp, err := http.Get(s.URL)
	if err != nil {
		log.Println(err.Error())
		return 500
	}

	return resp.StatusCode
}

func (s *Website) setSiteUp(dbConn *storage.Connection) {
	if !s.IsUp {
		s.endOutage(dbConn)
	}
	s.IsUp = true

	sql := "UPDATE sites SET last_checked = ?, is_up = ? WHERE url = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec(time.Now().Format("2006-01-02 15:04:05"), 1, s.URL)
	if err != nil {
		fmt.Println(err)
	}
	statement.Close()
}

func (s *Website) setSiteDown(dbConn *storage.Connection) {
	if s.IsUp {
		s.beginOutage(dbConn)
	}
	s.IsUp = false
	sql := "UPDATE sites SET last_checked = ?, is_up = ? WHERE url = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	_, err = statement.Exec(time.Now().Format("2006-01-02 15:04:05"), 0, s.URL)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Website) GetSiteID(dbConn *storage.Connection) int {
	var siteID int
	row, err := dbConn.DB.Query("SELECT id FROM sites WHERE url = '" + s.URL + "'")
	if err != nil {
		fmt.Println(err)
	}
	defer row.Close()
	for row.Next() {
		row.Scan(&siteID)
	}
	return siteID
}

func (s *Website) beginOutage(dbConn *storage.Connection) {
	siteID := s.GetSiteID(dbConn)
	sql := "insert into outages values (null, ?, ?, null);"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec(siteID, time.Now().Format("2006-01-02 15:04:05"))
}

func (s *Website) endOutage(dbConn *storage.Connection) {
	siteID := s.GetSiteID(dbConn)
	sql := "update outages set outage_end = ? where website_id = ? and outage_end is null"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec(time.Now().Format("2006-01-02 15:04:05"), siteID)
	statement.Close()
}

func FindWebsiteByID(id int, dbConn *storage.Connection) (Website, error) {
	s := Website{}
	row, err := dbConn.DB.Query("SELECT id, url, is_up FROM sites WHERE id = ?", id)
	if err != nil {
		fmt.Println(err)
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

// Destroy - delete a site from the DB and close down the monitoring routine
func (s *Website) Destroy(c *chan string, dbConn *storage.Connection) {
	sql := "UPDATE sites SET is_deleted = ? WHERE id = ?"
	statement, err := dbConn.DB.Prepare(sql)
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec(1, s.ID)
	statement.Close()
	go func() {
		*c <- s.URL
	}()
}
