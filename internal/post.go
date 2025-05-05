package internal

import (
	"time"
)

type Post struct {
	ID        string    `json:"id,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
