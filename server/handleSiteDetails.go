package server

import (
	"encoding/json"
	"net/http"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleSiteDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		siteURL := r.FormValue("url")
		if siteURL == "" {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		site, err := site.FindWebsiteByURL(siteURL, s.storage, s.logger)
		if err != nil {
			if err.Error() == "Site was not found when looping over DB records" {
				w.WriteHeader(404)
				w.Header().Add("Content-type", "application/json")
				w.Write([]byte("{\"error\":\"not found\"}"))
				return
			}
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}

		certData, _ := site.GetCertificateInfo()
		domainData, _ := site.GetDomainInfo()

		data := map[string]interface{}{
			"url":        site.URL,
			"domainInfo": domainData,
			"up":         site.IsUp,
			"outages":    site.Outages(s.storage),
			"uptime": map[string]interface{}{
				"days1":  site.CalcUptime(1, s.storage),
				"days7":  site.CalcUptime(7, s.storage),
				"days30": site.CalcUptime(30, s.storage),
				"days60": site.CalcUptime(60, s.storage),
				"days90": site.CalcUptime(90, s.storage),
			},
			"certInfo": certData,
		}
		dataJSON, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dataJSON)
	}
}
