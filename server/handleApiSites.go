package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleApiSites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.loadSites()
		bytes, _ := json.Marshal(s.sites)

		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	}
}
