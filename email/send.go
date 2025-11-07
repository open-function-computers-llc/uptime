package email

import (
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func Send(subject, message string, includeEmergencyPrefix bool) error {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAIL_FROM"))
	m.SetHeader("To", os.Getenv("EMAIL_TO"))

	if includeEmergencyPrefix {
		m.SetHeader("Subject", os.Getenv("EMERGENCY_EMAIL_PREFIX")+subject)
	} else {
		m.SetHeader("Subject", subject)
	}

	m.SetBody("text/html", message)

	smtpPortString := os.Getenv("SMTP_PORT")
	port, _ := strconv.Atoi(smtpPortString)
	d := gomail.NewDialer(
		os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASSWORD"),
	)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
