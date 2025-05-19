package usersworker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mrrobotisreal/journey-app-server-refresh/internal/cache"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
)

func HandleLogin(evt eventbus.Event) error {
	if evt.Type != eventbus.EventLogin {
		log.Println("Early exit because evt.Type is not eventbus.EventLogin")
		return nil
	}

	ctx := context.Background()
	fbID := evt.Firebase

	userID, username, _, err := db.Repo.GetUserByFirebaseLogin(ctx, fbID) // TODO: add apiKey
	if err != nil {
		log.Printf("Error fetching user by firebase ID: %v", err)
		return err
	}

	meta := map[string]any{
		"source": "backend",
	}
	blob, err := json.Marshal(meta)
	if err != nil {
		log.Printf("marshal metadata error: %v", err)
		return err
	}

	_, err = db.Repo.ExecContext(ctx,
		`INSERT INTO analytics_events (user_id, fb_id, event_type, meta_data, event_time)
		 VALUES (?, ?, 'login', ?, ?)`,
		userID, fbID, string(blob), time.Now().UTC())

	if err != nil {
		log.Printf("Error saving login analytics: %v", err)
	}

	cache.SaveUser(ctx, userID, fbID, username)

	log.Printf("[users-worker] processed login for user_id=%d", userID)
	return nil
}
