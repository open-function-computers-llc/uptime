package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleSiteDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		siteID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", 302)
		}
		site, err := site.FindWebsiteByID(siteID, s.storage, s.logger)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", 302)
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
