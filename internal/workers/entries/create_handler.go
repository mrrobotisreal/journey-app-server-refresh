package entriesworker

import (
	"context"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/cache"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/db"
	"github.com/mrrobotisreal/journey-app-server-refresh/internal/eventbus"
	models_entries_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/entries/create"
	"log"
)

func HandleCreateEntry(evt eventbus.Event) error {
	if evt.Type != eventbus.EventCreateEntry {
		log.Println("Early exit because evt.Type is not eventbus.EventCreateEntry")
		return nil
	}

	ctx := context.Background()
	entry := models_entries_create.CreateEntryRequest{
		ID:        evt.Payload["ID"].(string),
		UserID:    evt.Payload["userID"].(int64),
		Text:      evt.Payload["text"].(string),
		Locations: evt.Payload["locations"].([]models_entries_create.LocationData),
		Tags:      evt.Payload["tags"].([]models_entries_create.TagData),
		Images:    evt.Payload["images"].([]string),
	}
	cache.SaveEntry(ctx, entry)

	err := db.Repo.InsertEntry(ctx, entry)
	if err != nil {
		log.Printf("error inserting entry %s", entry.ID)
		return err
	}

	return nil
}
