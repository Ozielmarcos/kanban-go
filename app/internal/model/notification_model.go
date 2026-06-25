package model

type Notification struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
	Read    bool   `json:"read"`
}
