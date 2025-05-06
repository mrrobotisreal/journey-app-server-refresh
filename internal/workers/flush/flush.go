package flush

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/cache"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
)

func Handle(evt eventbus.Event) error {
	if evt.Type != eventbus.EventLogin {
		return nil
	}

	ctx := context.Background()

	key := fmt.Sprintf("user:%d", evt.UserID)
	vals, err := cache.RDB.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}

	blob, _ := json.Marshal(vals)
	_, err = db.Repo.ExecContext(ctx,
		`INSERT INTO analytics_events (user_id, fb_id, event_type, meta_data)
		 VALUES (?, ?, 'login', ?)`,
		evt.UserID, evt.Firebase, string(blob))
	return err
}
