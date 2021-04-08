package server

import (
	"net/http"
)

func (s *Server) logRequest(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// s.log(logLine(
		// 	time.Now().Format("2006-01-02 15:04:05"),
		// 	r.Method,
		// 	r.URL.Path,
		// 	r.UserAgent(),
		// ))
		h(w, r)
	}
}

func logLine(parts ...string) string {
	output := ""
	for _, part := range parts {
		output = output + part + " "
	}
	return output
}
