package cache

import (
	"context"
	"fmt"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	models_entries_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/entries/create"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // change localhost to redis when in Docker Compose network
})

func SaveUser(ctx context.Context, userID int64, fbID, name string) {
	key := fmt.Sprintf("user:%d", userID)
	RDB.HSet(ctx, key, map[string]any{
		"fb_id": fbID,
		"name":  name,
	}).Result()
	RDB.Expire(ctx, key, 24*time.Hour)
}

func TouchUser(ctx context.Context, userID int64) {
	RDB.Expire(ctx, fmt.Sprintf("user:%d", userID), 24*time.Hour)
}

func SaveEntry(ctx context.Context, entry models_entries_create.CreateEntryRequest) {
	key := fmt.Sprintf("entry:%s", entry.ID)
	RDB.HSet(ctx, key, map[string]any{
		"ID":        entry.ID,
		"userID":    entry.UserID,
		"text":      entry.Text,
		"locations": entry.Locations,
		"tags":      entry.Tags,
		"images":    entry.Images,
	}, 30*time.Minute).Result()
}

func HandleEntry(evt eventbus.Event) error {
	if evt.Type != eventbus.EventUpdateEntry && evt.Type != eventbus.EventDeleteEntry {
		return nil
	}
	id, ok := evt.Payload["entry_id"].(float64)
	if !ok {
		return fmt.Errorf("payload missing entry_id")
	}
	data, err := RDB.Get(context.Background(), fmt.Sprintf("entry:%d", int64(id))).Bytes()
	if err != nil {
		return err
	}

	return db.Repo.PersistEntry(context.Background(), int64(id), data)
}
