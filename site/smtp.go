package site

import (
	"errors"
	"os"
	"strconv"
)

func checkSMTPEnv() error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpPort := os.Getenv("SMTP_PORT")
	emailTo := os.Getenv("EMAIL_TO")
	emailFrom := os.Getenv("EMAIL_FROM")

	_, err := strconv.Atoi(smtpPort)
	if err != nil {
		return err
	}

	if smtpHost == "" || smtpUser == "" || smtpPass == "" || smtpPort == "" || emailTo == "" || emailFrom == "" {
		return errors.New("Missing required SMTP information")
	}

	return nil
}
