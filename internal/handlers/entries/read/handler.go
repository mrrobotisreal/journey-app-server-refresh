package handlers_entries_read

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
)

func ListEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("userID")
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "invalid userID", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "invalid page", http.StatusBadRequest)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	entries, err := db.Repo.ListEntries(r.Context(), userIDInt, pageInt, limitInt)
	if err != nil {
		http.Error(w, "failed to list entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
