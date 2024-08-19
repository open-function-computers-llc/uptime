package site

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func (s *Website) processWebhooks(t string) error {
	if t != "standard" && t != "emergency" {
		return errors.New("Invalid webhook type requested")
	}

	hooks, err := s.GetWebHooks(t)
	if err != nil {
		return err
	}

	if len(hooks) < 1 {
		return nil
	}

	for _, hook := range hooks {
		fmt.Println("about to " + hook.Verb + " " + hook.URL)

		client := &http.Client{}
		if hook.Verb == "POST" {
			body := bytes.NewBuffer([]byte(""))
			errorBody := ""
			req, err := http.NewRequest("POST", hook.URL, body)
			if err != nil {
				errorBody = err.Error()
			}

			resp, err := client.Do(req)
			if err != nil {
				errorBody = err.Error()
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				errorBody = err.Error()
			}

			m := gomail.NewMessage()
			m.SetHeader("From", os.Getenv("EMAIL_FROM"))
			m.SetHeader("To", os.Getenv("EMAIL_TO"))
			m.SetHeader("Subject", "Webhook Processed")
			if errorBody != "" {
				m.SetHeader("Subject", "Webhook Processed - ERROR")
			}
			m.SetBody("text/html", "<p>Webhook for "+hook.URL+" ran. Output:</p><p>"+errorBody+string(b)+"</p>")

			port := os.Getenv("SMTP_PORT")
			portInt, _ := strconv.Atoi(port)
			d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
				portInt,
				os.Getenv("SMTP_USER"),
				os.Getenv("SMTP_PASSWORD"))
			if err := d.DialAndSend(m); err != nil {
				s.Logger.Error(err)
			}
		}
	}

	return nil
}
