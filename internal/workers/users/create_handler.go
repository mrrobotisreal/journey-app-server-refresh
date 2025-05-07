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
	log.Printf("DEBUG got raw event: %+v", evt)
	if evt.Type != eventbus.EventCreateAccount {
		log.Println("Early exit because evt.Type is not eventbus.EventCreateAccount")
		return nil
	}

	log.Println("About to get context.Background()")
	ctx := context.Background()
	log.Println("Got context.Background()")

	log.Println("About to get evt.Payload[username].(string)")
	name, _ := evt.Payload["username"].(string)
	log.Println("Got evt.Payload[username].(string); about to cache.SaveUser")
	cache.SaveUser(ctx, evt.UserID, evt.Firebase, name)
	log.Println("Cached SaveUser")

	meta := map[string]any{
		"source": "backend",
	}
	blob, _ := json.Marshal(meta)
	log.Println("Marshaled meta")

	log.Println("About to db.Repo.ExecContext")
	_, err := db.Repo.ExecContext(ctx,
		`INSERT INTO analytics_events (user_id, fb_id, event_type, meta_data, event_time)
		 VALUES (?, ?, 'create_account', ?, ?)`,
		evt.UserID, evt.Firebase, string(blob), time.Now().UTC())
	log.Printf("Logging err after create_account event: %v", err)

	if err == nil {
		log.Printf("[users-worker] processed create_account for user_id=%d", evt.UserID)
	} else {
		log.Printf("Something went wrong trying to store create_account event analytics: %v", err)
	}
	return err
}
