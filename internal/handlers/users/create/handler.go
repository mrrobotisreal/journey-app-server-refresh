package handlers_users_create

import (
	"encoding/json"
	"firebase.google.com/go/v4/auth"
	"github.com/google/uuid"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/cache"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/firebase"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/glue/deps"
	models "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/users/create"
	"net/http"
	"time"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST required", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	fbUser, err := firebase.AuthClient.CreateUser(ctx, (&auth.UserToCreate{}).Email(req.Email).Password(req.Password).DisplayName(req.DisplayName))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, apiKey, err := db.Repo.InsertUser(ctx, fbUser.UID, req.Email, req.DisplayName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cache.SaveUser(ctx, userID, fbUser.UID, req.DisplayName)

	evt := eventbus.Event{
		ID:        uuid.New(),
		Type:      eventbus.EventCreateAccount,
		UserID:    userID,
		Firebase:  fbUser.UID,
		Payload:   map[string]any{"username": req.DisplayName},
		Timestamp: time.Now().UTC(),
	}
	_ = deps.Bus.Publish(ctx, "users", evt)

	json.NewEncoder(w).Encode(models.CreateAccountResponse{userID, apiKey})
}
