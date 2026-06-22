package model

import "time"

type taskStatus string

const (
	ToDo       taskStatus = "to_do"
	InProgress taskStatus = "in_progress"
	Done       taskStatus = "done"
	Review     taskStatus = "review"
	Blocked    taskStatus = "blocked"
)

type taskPriority string

const (
	Low    taskPriority = "baixa"
	Medium taskPriority = "media"
	High   taskPriority = "alta"
)

type TimeEntry struct {
	ID        string `json:"id"`         //UUID
	TaskId    string `json:"task_id"`    //UUID
	StartTime string `json:"start_time"` // ISO string
	EndTime   string `json:"end_time"`   // ISO string
	Day       string `json:"day"`        // YYYY-MM-DD
}

type Task struct {
	ID                string       `json:"id"`
	StoryId           string       `json:"story_id"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	Effort            float32      `json:"effort"`
	Status            taskStatus   `json:"status"`
	Priority          taskPriority `json:"priority"`
	CreatedAt         time.Time    `json:"created_at"` // ISO string
	TimeEntries       []TimeEntry  `json:"time_entries,omitempty"`
	IsTimerRunning    bool         `json:"is_timer_running,omitempty"`
	CurrentTimerStart *time.Time   `json:"current_timer_start,omitempty"` // ISO string
}
