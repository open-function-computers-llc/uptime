package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleRestoreSite() http.HandlerFunc {
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
		s.log(site)
		site.Restore(s.storage, s.logger)
		site.Monitor(s.shutdownChannel)

		http.Redirect(w, r, "/", 302)
	}
}
