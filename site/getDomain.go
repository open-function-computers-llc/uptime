package site

import (
	"strings"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

func (s *Website) GetDomain() string {
	domain := strings.Replace(s.URL, "http://www.", "", 1)
	domain = strings.Replace(domain, "https://www.", "", 1)
	domain = strings.Replace(domain, "http://", "", 1)
	domain = strings.Replace(domain, "https://", "", 1)

	if domain[len(domain)-1:] == "/" {
		domain = domain[:len(domain)-1]
	}
	return domain
}

type WhoisInfo struct {
	Status          []string `json:"status"`
	RegisteredDate  string   `json:"registered"`
	ExpiresDate     string   `json:"expires"`
	RegisteredName  string   `json:"registrant"`
	RegisteredEmail string   `json:"registrantEmail"`
	RegisteredAt    string   `json:"registrar"`
	Error           string   `json:"error"`
}

func (s *Website) GetDomainInfo() (WhoisInfo, error) {
	output := WhoisInfo{}
	res, err := whois.Whois(s.GetDomain())
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
