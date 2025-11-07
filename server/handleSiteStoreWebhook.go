package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteStoreWebhook() http.HandlerFunc {
	type incomingPayload struct {
		SiteID   int    `json:"siteID"`
		Name     string `json:"name"`
		Url      string `json:"url"`
		Verb     string `json:"verb"`
		HookType string `json:"hookType"`
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

		err = storage.StoreWebhook(payload.Name, payload.Url, payload.Verb, strings.ToLower(payload.HookType), payload.SiteID, s.storage)

		http.Redirect(w, r, "/webhooks/"+strconv.Itoa(payload.SiteID), http.StatusFound)
	}
}
