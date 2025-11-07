package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleSiteUpdate() http.HandlerFunc {
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

		siteID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Redirect(w, r, "/?error=not valid", http.StatusFound)
			return
		}

		s.sites[siteID].URL = payload.URL
		s.sites[siteID].Meta = payload.Meta
		storage.UpdateSite(siteID, payload.URL, payload.Meta, s.storage)

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
