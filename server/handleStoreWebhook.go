package server

import (
	"net/http"

	"github.com/open-function-computers-llc/uptime/webhook"
)

func (s *Server) handleSiteStoreWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		name := r.Form.Get("new-hook-name")
		url := r.Form.Get("new-hook-url")
		method := r.Form.Get("new-hook-method")
		hookType := r.Form.Get("new-hook-type")
		siteID := r.Form.Get("site-id")

		err := webhook.StoreWebhook(s.storage.DB, name, url, method, hookType, siteID)
		if err != nil {
			// TODO: handle this error?
			s.logger.Error(err)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
