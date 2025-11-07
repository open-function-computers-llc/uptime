package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteDeleteWebhook() http.HandlerFunc {
	type incomingPayload struct {
		SiteID    int `json:"siteID"`
		WebhookID int `json:"webhookID"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		var payload incomingPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err = storage.DeleteWebhook(payload.WebhookID, s.storage)
		if err != nil {
			http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/webhooks/"+strconv.Itoa(payload.SiteID), http.StatusSeeOther)
	}
}
