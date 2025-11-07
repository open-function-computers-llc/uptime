package server

import (
	"net/http"
	"strconv"
)

func (s *server) handleSitePurgeCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		for id, site := range s.sites {
			if id == siteID {
				s.inertiaManager.Render(w, r, "PurgeCheck", map[string]any{
					"site": site,
				})
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
