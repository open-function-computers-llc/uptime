package server

import (
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSitePurge() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		for id, site := range s.sites {
			if id == siteID {
				storage.StopMonitoringSite(site, s.storage.Logger)
			}
		}

		err = storage.PurgeSite(siteID, s.storage)
		if err != nil {
			s.log(err.Error())
		}
		delete(s.sites, siteID)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
