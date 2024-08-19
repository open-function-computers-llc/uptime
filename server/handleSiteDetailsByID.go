package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

// Depreciated. Check out the new handler s.handleSiteDetails()
func (s *Server) handleSiteDetailsByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}
		site, err := site.FindWebsiteByID(siteID, s.storage, s.logger)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data := map[string]interface{}{
			"url":     site.URL,
			"up":      site.IsUp,
			"outages": site.Outages(s.storage),
			"uptime": map[string]interface{}{
				"days1":  site.CalcUptime(1, s.storage),
				"days7":  site.CalcUptime(7, s.storage),
				"days30": site.CalcUptime(30, s.storage),
				"days60": site.CalcUptime(60, s.storage),
				"days90": site.CalcUptime(90, s.storage),
			},
		}
		dataJSON, _ := json.Marshal(data)
		w.Write(dataJSON)
	}
}
