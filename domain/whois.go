package domain

import (
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

type WhoisInfo struct {
	Status          []string `json:"status"`
	RegisteredDate  string   `json:"registered"`
	ExpiresDate     string   `json:"expires"`
	RegisteredName  string   `json:"registrant"`
	RegisteredEmail string   `json:"registrantEmail"`
	RegisteredAt    string   `json:"registrar"`
	Error           string   `json:"error"`
}

func GetDomainInfo(url string) (WhoisInfo, error) {
	output := WhoisInfo{}
	res, err := whois.Whois(GetDomainFromURL(url))
	if err != nil {
		return output, err
	}

	parsed, err := whoisparser.Parse(res)
	if err != nil {
		return output, err
	}

	output.Status = parsed.Domain.Status
	output.RegisteredName = parsed.Registrant.Name
	output.RegisteredEmail = parsed.Registrant.Email
	output.RegisteredAt = parsed.Registrar.Name
	output.RegisteredDate = parsed.Domain.CreatedDate
	output.ExpiresDate = parsed.Domain.ExpirationDate

	return output, nil
}
