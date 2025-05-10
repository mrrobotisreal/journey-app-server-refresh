package models_entries_create

import "time"

type LocationData struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	DisplayName string  `json:"displayName"`
}

type TagData struct {
	Key   string  `json:"key"`
	Value *string `json:"value"`
}

type CreateEntryRequest struct {
	ID        string         `json:"ID"`
	UserID    int64          `json:"userID"`
	Text      string         `json:"text"`
	Locations []LocationData `json:"locations"`
	Tags      []TagData      `json:"tags"`
	Images    []string       `json:"images"`
}

type Entry struct {
	ID        string         `json:"ID"`
	UserID    int64          `json:"userID"`
	Text      string         `json:"text"`
	Locations []LocationData `json:"locations"`
	Tags      []TagData      `json:"tags"`
	Images    []string       `json:"images"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}
