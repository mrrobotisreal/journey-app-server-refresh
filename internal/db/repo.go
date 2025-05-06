package db

import (
	"context"
	"database/sql"

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

/******** Entries ********/

func (r *Repository) PersistEntry(ctx context.Context, entryID int64, data []byte) error {
	// TODO: update this accordingly later
	_, err := r.ExecContext(ctx,
		`UPDATE entries SET text = ? WHERE entry_id = ?`,
		string(data), entryID)
	return err
}
