package eventbus

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventLogin         EventType = "login"
	EventCreateEntry   EventType = "create_entry"
	EventUpdateEntry   EventType = "update_entry"
	EventDeleteEntry   EventType = "delete_entry"
	EventReadEntry     EventType = "read_entry"
	EventCreateAccount EventType = "create_account"
	EventDeleteAccount EventType = "delete_account"
)

type Event struct {
	ID        uuid.UUID      `json:"id"`
	Type      EventType      `json:"type"`
	UserID    int64          `json:"user_id"`
	Firebase  string         `json:"firebase"`
	Payload   map[string]any `json:"payload,omitempty"`
	Timestamp time.Time      `json:"ts"`
}

func New(t EventType, uid int64, p map[string]any) Event {
	return Event{
		ID:        uuid.New(),
		Type:      t,
		UserID:    uid,
		Payload:   p,
		Timestamp: time.Now().UTC(),
	}
}
