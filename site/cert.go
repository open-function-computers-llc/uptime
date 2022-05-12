package site

import (
	"crypto/tls"
	"fmt"
	"time"
)

type CertificateInfo struct {
	Valid   bool      `json:"valid"`
	Expires time.Time `json:"expires"`
	Info    string    `json:"error"`
	Names   []string  `json:"names"`
}

func (s *Website) GetCertificateInfo() (CertificateInfo, error) {
	info := CertificateInfo{
		Names: []string{},
	}
	conn, err := tls.Dial("tcp", s.GetDomain()+":443", nil)
	if err != nil {
		fmt.Println(err)
		return info, err
	}
	defer conn.Close()

	err = conn.VerifyHostname(s.GetAddressWithoutProtocol())
	if err != nil {
		info.Valid = false
		info.Info = "Hostname doesn't match site URL: " + err.Error()
		return info, err
	}

	expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
	for _, chain := range conn.ConnectionState().VerifiedChains {
		for _, cert := range chain {
			for _, name := range cert.DNSNames {
				if name == "" {
					continue
				}

				exists := false
				for _, existingName := range info.Names {
					if name == existingName {
						exists = true
						continue
					}
				}
				if !exists {
					info.Names = append(info.Names, name)
				}
			}
		}
	}
	info.Valid = true
	info.Expires = expiry
	return info, nil
}
