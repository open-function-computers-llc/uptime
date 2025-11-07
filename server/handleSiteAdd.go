package server

import (
	"net/http"
)

func (s *server) handleSiteAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.inertiaManager.Render(w, r, "Create", map[string]any{})
	}
}
