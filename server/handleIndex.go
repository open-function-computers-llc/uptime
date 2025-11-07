package server

import (
	"net/http"
)

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.loadSites()
		s.inertiaManager.Render(w, r, "Index", map[string]any{
			"sites": s.sites,
		})
	}
}
