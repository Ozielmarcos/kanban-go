package model

type TimeEntry struct {
	ID        string  `json:"id"`      //UUID
	TaskId    string  `json:"task_id"` //UUID
	TaskTitle string  `json:"task_title"`
	StartTime string  `json:"start_time"` // ISO string
	EndTime   string  `json:"end_time"`   // ISO string
	Day       string  `json:"day"`        // YYYY-MM-DD
	Duration  float64 `json:"duration,omitempty"`
}
