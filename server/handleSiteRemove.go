package server

import (
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		site := s.sites[siteID]
		storage.StopMonitoringSite(site, s.storage.Logger)
		storage.DeleteSite(s.storage, site.ID)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
