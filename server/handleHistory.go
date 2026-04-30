package server

import (
	"net/http"
)

func (s *server) handleHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		countErrMessage := ""
		secErrMessage := ""

		outages, err := s.outagesByDay()
		if err != nil {
			countErrMessage = err.Error()
		}
		outageDurations, err := s.outagesDurationsByDay()
		if err != nil {
			secErrMessage = err.Error()
		}

		s.inertiaManager.Render(w, r, "History", map[string]any{
			"outages":         outages,
			"outageDurations": outageDurations,
			"countErrMessage": countErrMessage,
			"secErrMessage":   secErrMessage,
		})
	}
}
