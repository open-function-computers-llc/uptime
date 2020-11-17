package server

import "net/http"

func (s *Server) logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.log(r.Method+": ", r)
		h(w, r)
	}
}
