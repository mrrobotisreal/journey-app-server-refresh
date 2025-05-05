package models_entries_delete

import "time"

type DeleteImageRequest struct {
	Username      string    `json:"username"`
	UserID        int64     `json:"userId"`
	Timestamp     time.Time `json:"timestamp"`
	EntryID       int64     `json:"entryId"`
	Images        []string  `json:"images"`
	ImageToDelete string    `json:"imageToDelete"`
}

type DeleteImageResponse struct {
	Success bool `json:"success"`
}
