package repository

import (
	"database/sql"
	"errors"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/pkg/database"
)

func GetTasks() []model.Task {
	_, err := database.DB.Query("SELECT * FROM tasks")
	if err != nil {
		return []model.Task{}
	}
	return []model.Task{}
}

func GetTasksByUser(storyId string) ([]model.Task, error) {
	rows, err := database.DB.Query("SELECT * FROM tasks WHERE story_id = $1", storyId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.StoryId, &task.Title, &task.Description, &task.Effort, &task.Status, &task.Priority, &task.CreatedAt, &task.IsTimerRunning, &task.CurrentTimerStart)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskById(id string, storyId string) (model.Task, error) {
	var task model.Task

	err := database.DB.QueryRow(
		"SELECT * FROM tasks WHERE id = $1 AND story_id = $2",
		id, storyId,
	).Scan(&task.ID, &task.StoryId, &task.Title, &task.Description, &task.Effort, &task.Status, &task.Priority, &task.CreatedAt, &task.IsTimerRunning, &task.CurrentTimerStart)

	return task, err
}

func CreateTask(task model.Task) (model.Task, error) {
	err := database.DB.QueryRow(
		`INSERT INTO tasks ( story_id, title, description, effort, status, priority) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING 
				id,
				story_id, 
				title, 
				description, 
				effort, 
				status, 
				priority, 
				created_at,
				is_timer_running,
    			current_timer_start`,
		task.StoryId, task.Title, task.Description, task.Effort, task.Status, task.Priority,
	).Scan(
		&task.ID,
		&task.StoryId,
		&task.Title,
		&task.Description,
		&task.Effort,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.IsTimerRunning,
		&task.CurrentTimerStart,
	)

	return task, err
}

func UpdateTask(id string, task model.Task) error {
	sql := "UPDATE tasks SET title = $1, description = $2, effort = $3, status = $4, priority = $5 WHERE id = $6 AND story_id = $7"
	_, err := database.DB.Exec(sql, task.Title, task.Description, task.Effort, task.Status, task.Priority, id, task.StoryId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateTaskStatus(id string, status model.TaskStatus) error {
	sql := "UPDATE tasks SET status=$1 WHERE id=$2"

	_, err := database.DB.Exec(sql, status, id)

	return err
}

func StarTimer(taskId string) error {
	sql := "UPDATE tasks SET is_timer_running = true, current_timer_start = NOW() where id = $1"

	_, err := database.DB.Exec(sql, taskId)
	if err != nil {
		return err
	}

	return nil
}

func PauseTimer(taskId string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var startTimer sql.NullTime

	err = tx.QueryRow(`
		SELECT current_timer_start FROM tasks
		WHERE id = $1
		AND is_timer_running = true
		`, taskId,
	).Scan(&startTimer)

	if err != nil {
		return err
	}

	if !startTimer.Valid {
		return errors.New("Timer não esta em execução!")
	}

	_, err = tx.Exec(
		`INSERT INTO time_entries (
			task_id,
			start_time,
			end_time,
			day
		) VALUES ($1, $2, NOW(), CURRENT_DATE);`,
		taskId, startTimer.Time,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		`UPDATE tasks
			SET is_timer_running = false, 
			current_timer_start = null
		WHERE id = $1;	
		`, taskId)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func ResumeTimer(taskId string) error {
	_, err := database.DB.Exec(`
		UPDATE tasks
		SET is_timer_running = true, current_timer_start = NOW()
		WHERE id = $1;
	`, taskId)

	if err != nil {
		return err
	}

	return nil
}
