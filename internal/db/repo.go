package db

import (
	"context"
	"database/sql"
	"log"

	models_entries_create "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/entries/create"

	"github.com/google/uuid"
)

type Repository struct{ *sql.DB }

var Repo *Repository

/******** USERS ********/

func (r *Repository) InsertUser(ctx context.Context, fbID, email, username string) (int64, string, error) {
	apiKey := uuid.New().String()
	res, err := r.ExecContext(ctx, `
		INSERT INTO users (fb_id, email, username, api_key)
		VALUES (?, ?, ?, ?)`,
		fbID, email, username, apiKey)
	if err != nil {
		return 0, "", err
	}
	uid, _ := res.LastInsertId()
	return uid, apiKey, nil
}

func (r *Repository) GetUserByFirebase(ctx context.Context, fbID string) (int64, error) {
	var id int64
	err := r.QueryRowContext(ctx,
		`SELECT user_id FROM users WHERE fb_id = ?`, fbID).Scan(&id)
	return id, err
}

func (r *Repository) GetUserByFirebaseLogin(ctx context.Context, fbID string) (int64, string, string, error) {
	var userID int64
	var username, apiKey string
	err := r.QueryRowContext(ctx,
		`SELECT user_id, username, api_key FROM users WHERE fb_id = ?`, fbID).Scan(&userID, &username, &apiKey)

	if err == nil {
		_, _ = r.ExecContext(ctx,
			`UPDATE users SET api_key_last_used_at = CURRENT_TIMESTAMP WHERE user_id = ?`, userID)
	}

	return userID, username, apiKey, err
}

/******** Entries ********/

func (r *Repository) PersistEntry(ctx context.Context, entryID int64, data []byte) error {
	// TODO: update this accordingly later
	_, err := r.ExecContext(ctx,
		`UPDATE entries SET text = ? WHERE entry_id = ?`,
		string(data), entryID)
	return err
}

func (r *Repository) InsertEntry(ctx context.Context, entry models_entries_create.CreateEntryRequest) error {
	_, err := r.ExecContext(ctx, `
		INSERT INTO entries (entry_id, user_id, text)
		VALUES (?, ?, ?)
	`, entry.ID, entry.UserID, entry.Text)
	if err != nil {
		return err
	}

	if len(entry.Locations) > 0 {
		for _, loc := range entry.Locations {
			_, err = r.ExecContext(ctx, `
				INSERT INTO entry_locations (entry_id, latitude, longitude, display_name)
				VALUES (?, ?, ?, ?)
			`, entry.ID, loc.Latitude, loc.Longitude, loc.DisplayName)
			if err != nil {
				// TODO: handle this unhappy path better
				log.Printf("error inserting location for entry %s, lat = %f, lng = %f, dName = %s", entry.ID, loc.Latitude, loc.Longitude, loc.DisplayName)
			}
		}
	}

	if len(entry.Tags) > 0 {
		for _, tag := range entry.Tags {
			_, err = r.ExecContext(ctx, `
				INSERT INTO entry_tags (entry_id, tag_key, tag_value)
				VALUES (?, ?, ?)
			`, entry.ID, tag.Key, tag.Value)
			if err != nil {
				// TODO: handle this unhappy path better
				log.Printf("error inserting tag for entry %s, key = %s, val = %v", entry.ID, tag.Key, tag.Value)
			}
		}
	}

	if len(entry.Images) > 0 {
		for _, img := range entry.Images {
			_, err = r.ExecContext(ctx, `
				INSERT INTO entry_images (entry_id, image_url)
				VALUES (?, ?)
			`, entry.ID, img)
			if err != nil {
				// TODO: handle this unhappy path better
				log.Printf("error inserting image for entry %s, url = %s", entry.ID, img)
			}
		}
	}

	return nil
}

func (r *Repository) ListEntries(ctx context.Context, userID, page, limit int) ([]models_entries_create.Entry, error) {
	offset := (page - 1) * limit
	rows, err := r.QueryContext(ctx, `
		SELECT entry_id, user_id, text, created_at, updated_at
		FROM entries
		ORDER BY created_at DESC
		WHERE user_id = ?
		LIMIT ? OFFSET ?
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []models_entries_create.Entry{}
	for rows.Next() {
		var entry models_entries_create.Entry
		err = rows.Scan(&entry.ID, &entry.UserID, &entry.Text, &entry.CreatedAt, &entry.UpdatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
