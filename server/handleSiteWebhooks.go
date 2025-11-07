package server

import (
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteWebhooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		for id, site := range s.sites {
			if id == siteID {
				webhooks, err := storage.SiteWebhooks(site.ID, s.storage)
				if err != nil {
					s.log(err.Error())
					http.Redirect(w, r, "/?error=not valid", http.StatusFound)
					return
				}

				s.inertiaManager.Render(w, r, "Webhooks", map[string]any{
					"webhooks": webhooks,
					"site":     site,
				})
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
