package handlers_entries_create

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/glue/deps"
	models "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/entries/create"
	"net/http"
	"time"
)

func CreateEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateEntryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	evt := eventbus.Event{
		ID:     uuid.New(),
		Type:   eventbus.EventCreateEntry,
		UserID: req.UserID,
		Payload: map[string]any{
			"ID":        req.ID,
			"text":      req.Text,
			"locations": req.Locations,
			"tags":      req.Tags,
			"images":    req.Images,
		},
		Timestamp: time.Now().UTC(),
	}
	_ = deps.Bus.Publish(ctx, "entries", evt)
}
