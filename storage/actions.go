package storage

import (
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/uptime/email"
	"github.com/open-function-computers-llc/uptime/models"
)

func processStandardActions(secondsDown int, s *models.Site, webhooks []*models.Webhook) error {
	if s.StandardWarningSent {
		return nil
	}

	err := email.Send(s.URL+" is down!", buildHTMLMessage(s.URL, s.Meta, secondsDown), false)
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

	err := email.Send(s.URL+" IS STILL DOWN!", buildHTMLMessage(s.URL, s.Meta, secondsDown), true)
	if err != nil {
		return err
	}

	for _, wh := range webhooks {
		if wh.Type == models.HooktypeStandard {
			continue
		}
		wh.Process()
	}

	s.EmergencyWarningSent = true

	return nil
}

func buildHTMLMessage(url, meta string, secondsDown int) string {
	return `
		<h1>` + url + ` is down!</h1>
		<p>It has been down for at least ` + strconv.Itoa(secondsDown) + ` seconds. Better go check things out...</p>
		<p><strong>Meta info:</strong><br />` + strings.ReplaceAll(meta, "\n", "<br />") + `</p>
	`
}
