package handlers_users_login

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/firebase"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/glue/deps"
	models "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/users/login"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	fbResult, err := firebase.SignInWithEmail(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	ctx := r.Context()
	evt := eventbus.Event{
		ID:        uuid.New(),
		Type:      eventbus.EventLogin,
		Firebase:  fbResult.LocalID,
		Payload:   map[string]any{"email": req.Email},
		Timestamp: time.Now().UTC(),
	}

	err = deps.Bus.Publish(ctx, "users", evt)
	if err != nil {
		http.Error(w, "failed to publish event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"firebase": fbResult.LocalID,
		"email":    fbResult.Email,
		"idToken":  fbResult.IDToken,
	})
}
