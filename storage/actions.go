package storage

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/uptime/email"
	"github.com/open-function-computers-llc/uptime/models"
)

func processStandardActions(secondsDown int, s *models.Site, webhooks []*models.Webhook) error {
	if s.StandardWarningSent {
		return nil
	}

	err := email.Send(s.URL+" is down!", buildHTMLDownMessage(s.URL, s.Meta, secondsDown), false)
	if err != nil {
		return err
	}

	for _, wh := range webhooks {
		if wh.Type == models.HooktypeEmergency {
			continue
		}
		wh.Process()
	}

	s.StandardWarningSent = true

	return nil
}

func processEmergencyActions(secondsDown int, s *models.Site, webhooks []*models.Webhook) error {
	if s.EmergencyWarningSent {
		return nil
	}

	// this either worked or didn't, so whatever
	go sendGotifyHTTPRequest(s.URL)

	err := email.Send(s.URL+" IS STILL DOWN!", buildHTMLDownMessage(s.URL, s.Meta, secondsDown), true)
	if err != nil {
		return err
	}

	for _, wh := range webhooks {
		if wh.Type == models.HooktypeStandard {
			continue
		}
		go wh.Process()
	}

	s.EmergencyWarningSent = true

	return nil
}

func sendGotifyHTTPRequest(domain string) {
	sendGotifyRequest := os.Getenv("GOTIFY_NOTIFICATION")
	if sendGotifyRequest != "true" {
		return
	}

	data := url.Values{
		"message":  {os.Getenv("EMERGENCY_EMAIL_PREFIX") + domain + " is Down!"},
		"priority": {"8"},
	}

	http.PostForm(os.Getenv("GOTIFY_HOST")+"message?token="+os.Getenv("GOTIFY_TOKEN"), data)
}

func buildHTMLDownMessage(url, meta string, secondsDown int) string {
	return `
	<h1>` + url + ` is down!</h1>
	<p>It has been down for at least ` + strconv.Itoa(secondsDown) + ` seconds. Better go check things out...</p>
	<p><strong>Meta info:</strong><br />` + strings.ReplaceAll(meta, "\n", "<br />") + `</p>
	`
}

func buildHTMLUpMessage(url, meta string, secondsDown int) string {
	return `
	<h1>` + url + ` is back up!</h1>
	<p>It has was down for ` + strconv.Itoa(secondsDown) + ` seconds.</p>
	<p><strong>Meta info:</strong><br />` + strings.ReplaceAll(meta, "\n", "<br />") + `</p>
	`
}
