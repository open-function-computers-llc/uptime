package server

import (
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handlePurgeSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		site, err := site.FindWebsiteByID(siteID, s.storage, s.logger)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		s.log(site)
		site.Purge(s.storage, s.logger)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
