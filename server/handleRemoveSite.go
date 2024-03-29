package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleRemoveSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		siteID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		site, err := site.FindWebsiteByID(siteID, s.storage, s.logger)
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
		}
		s.log(site)
		site.Destroy(s.shutdownChannel, s.storage, s.logger)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
