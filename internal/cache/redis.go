package cache

import (
	"context"
	"fmt"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	"time"

	"github.com/redis/go-redis/v9"
)

var RDB = redis.NewClient(&redis.Options{
	Addr: "redis:6379",
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

func SaveEntry(ctx context.Context, entryID int64, data []byte) {
	key := fmt.Sprintf("entry:%d", entryID)
	RDB.Set(ctx, key, data, 15*time.Minute)
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
