package site

import (
	"net/http"
	"os"
	"strings"
	"time"
)

func (s *Website) getStatusCode() int {
	if s.URL == "" {
		return 404
	}
	timeoutSeconds := 10

	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	req, err := http.NewRequest("GET", s.URL, nil)
	// TODO: this error checking is all wrong... oops!
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
