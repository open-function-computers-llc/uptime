package domain

import "strings"

func GetDomainFromURL(url string) string {
	domain := strings.Replace(url, "http://", "", 1)
	domain = strings.Replace(domain, "https://", "", 1)

	if domain[len(domain)-1:] == "/" {
		domain = domain[:len(domain)-1]
	}

	parts := strings.Split(domain, ".")
	if len(parts) > 2 {
		return strings.Join(parts[len(parts)-2:], ".")
	}

	return domain
}
