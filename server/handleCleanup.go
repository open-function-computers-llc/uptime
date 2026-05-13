package server

import (
	"encoding/json"
	"net/http"

	"github.com/open-function-computers-llc/uptime/storage"
)

func (s *server) handleCleanup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Check for active outages
		count, err := storage.GetCurrentOutageCount(s.storage)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Failed to check active outages: " + err.Error(),
			})
			return
		}

		// 2. If active outages exist, refuse cleanup
		if count > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  "error",
				"message": "Cannot clean up. There are currently " + string(rune(count)) + " active outage(s).",
				"active":  count,
			})
			return
		}

		// 3. Run cleanup
		err = storage.CleanupOutages(s.storage)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Cleanup failed: " + err.Error(),
			})
			return
		}

		// 4. Success
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "success",
			"message": "Outages cleaned up and IDs reset successfully.",
		})
	}
}
