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

func HandleCreateAccount(evt eventbus.Event) error {
	if evt.Type != eventbus.EventCreateAccount {
		return nil
	}

	ctx := context.Background()

	name, _ := evt.Payload["username"].(string)
	cache.SaveUser(ctx, evt.UserID, evt.Firebase, name)

	meta := map[string]any{
		"source": "backend",
	}
	blob, _ := json.Marshal(meta)

	_, err := db.Repo.ExecContext(ctx,
		`INSERT INTO analytics_events (user_id, fb_id, event_type, meta_data, event_time)
		 VALUES (?, ?, 'create_account', ?, ?)`,
		evt.UserID, evt.Firebase, string(blob), time.Now().UTC())

	if err == nil {
		log.Printf("[users-worker] processed create_account for user_id=%d", evt.UserID)
	} else {
		log.Printf("Something went wrong trying to store create_account event analytics: %v", err)
	}
	return err
}
