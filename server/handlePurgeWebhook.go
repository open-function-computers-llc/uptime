package server

import (
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/webhook"
)

func (s *Server) handlePurgeWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webhookID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		err = webhook.Delete(s.storage.DB, webhookID)
		if err != nil {
			s.logger.Error(err)
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
