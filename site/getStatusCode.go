package site

import (
	"net/http"
	"os"
	"strings"
	"time"
)

func (s *Website) getStatusCodeAndErrorMessage() (int, string) {
	if s.URL == "" {
		return 404, "404: URL is blank"
	}

	timeoutSeconds := 10

	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	req, err := http.NewRequest("GET", s.URL, nil)

	if err != nil {
		if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
			// s.Logger.Error(s.URL + " took too long to respond, and the URL was different!")
			return timeoutSeconds, err.Error()
		}
		return timeoutSeconds, err.Error()
	}
	req.Header.Set("User-Agent", "OFC_Uptime_Bot-version:"+os.Getenv("VERSION"))

	resp, err := client.Do(req)
	if err != nil {
		if strings.HasSuffix(err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
			// s.Logger.Error(s.URL + " took too long to respond!")
			return timeoutSeconds, err.Error()
		}
		return 500, err.Error()
	}
	defer resp.Body.Close()

	return resp.StatusCode, ""
}
