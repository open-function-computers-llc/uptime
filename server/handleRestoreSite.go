package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleRestoreSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server", s.storage)
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
		site.Restore(s.storage, s.logger)
		site.Monitor(s.shutdownChannel)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
