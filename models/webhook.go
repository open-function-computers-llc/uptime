package models

import (
	"bytes"
	"io"
	"net/http"

	"github.com/open-function-computers-llc/uptime/email"
)

type Webhook struct {
	ID        int
	WebsiteID int
	Name      string
	URL       string
	Verb      string
	Type      Hooktype
}

type Hooktype int

var HooktypeStandard Hooktype = 1
var HooktypeEmergency Hooktype = 2

func (wh *Webhook) Process() error {
	body := bytes.NewBuffer([]byte(""))

	req, err := http.NewRequest(wh.Verb, wh.URL, body)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return email.Send("Webhook Processed - "+wh.Name, "Webhook processed ("+wh.URL+"). Output:<br /><br />"+string(b), false)
}
