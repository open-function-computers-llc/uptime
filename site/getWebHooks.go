package site

import (
	"errors"
	"fmt"

	"github.com/open-function-computers-llc/uptime/webhook"
)

func (s *Website) GetWebHooks(t string) ([]webhook.Webhook, error) {
	hooks := []webhook.Webhook{}

	if t != "standard" && t != "emergency" {
		return hooks, errors.New("Invalid webhook type requested")
	}

	hookType := 1 // standard
	if t == "emergency" {
		hookType = 2
	}

	row, err := s.DB.DB.Query("SELECT id, website_id, hook_name, hook_url, hook_verb, hook_type FROM webhooks WHERE website_id = ? and hook_type = ?", s.ID, hookType)
	if err != nil {
		fmt.Println(err.Error())
		return hooks, err
	}
	defer row.Close()

	for row.Next() {
		var id int
		var wsID int
		var hname string
		var hurl string
		var hverb string
		var hooktype int

		row.Scan(&id, &wsID, &hname, &hurl, &hverb, &hooktype)

		hook := webhook.Webhook{
			ID:        id,
			WebsiteID: wsID,
			Name:      hname,
			URL:       hurl,
			Verb:      hverb,
			Type:      hooktype,
		}
		hooks = append(hooks, hook)
	}
	return hooks, nil
}
