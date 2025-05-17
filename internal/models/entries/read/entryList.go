package models_entries_read

import models "github.com/mrrobotisreal/journey-app-server-refresh/internal/models/entries/create"

type EntryList struct {
	Entries []models.Entry `json:"entries"`
	Total   int64          `json:"total"`
	Page    int64          `json:"page"`
	Limit   int64          `json:"limit"`
}
