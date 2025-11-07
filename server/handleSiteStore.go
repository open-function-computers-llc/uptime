package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteStore() http.HandlerFunc {
	type incomingPayload struct {
		URL  string `json:"url"`
		Meta string `json:"meta"`
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

		site, _ := storage.CreateSite(payload.URL, payload.Meta, s.storage)
		s.log(site)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
