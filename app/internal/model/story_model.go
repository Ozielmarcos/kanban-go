package model

import "time"

type Story struct {
	ID          string    `json:"id"`
	UserId      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
