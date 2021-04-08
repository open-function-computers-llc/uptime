package server

import (
	"net/http"

	"github.com/open-function-computers-llc/uptime/site"
)

func (s *Server) handleStoreSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		url := r.Form.Get("url")
		site := site.Create(url, s.storage, s.logger)
		site.Monitor(s.shutdownChannel)
		http.Redirect(w, r, "/", 302)
	}
}
