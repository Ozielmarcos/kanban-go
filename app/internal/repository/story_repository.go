package repository

import (
	"fmt"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/pkg/database"
)

func CreateStory(story model.Story) (model.Story, error) {
	err := database.DB.QueryRow(
		"INSERT INTO stories ( user_id, title, description) VALUES ($1, $2, $3) RETURNING id, user_id, title, description, created_at",
		story.UserId, story.Title, story.Description,
	).Scan(
		&story.ID,
		&story.UserId,
		&story.Title,
		&story.Description,
		&story.CreatedAt,
	)

	return story, err
}

func GetStoriesByUser(userId string) ([]model.Story, error) {
	// 1. Recomendado especificar as colunas na ordem exata do Scan
	query := "SELECT id, user_id, title, description FROM stories WHERE user_id = $1"
	rows, err := database.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Inicializa como slice vazio em vez de nil para o JSON vir como [] e não null
	stories := []model.Story{}

	for rows.Next() {
		var story model.Story
		err := rows.Scan(&story.ID, &story.UserId, &story.Title, &story.Description)
		if err != nil {
			// REMOVA O 'continue' TEMPORARIAMENTE PARA LOGAR O ERRO REAL
			fmt.Printf("ERRO NO SCAN DO BANCO: %v\n", err)
			return nil, err
		}
		stories = append(stories, story)
	}

	// 2. Verifica se o loop terminou por algum erro oculto do banco
	if err = rows.Err(); err != nil {
		fmt.Printf("ERRO NO CURSOR ROWS: %v\n", err)
		return nil, err
	}

	return stories, nil
}

func UpdateStory(id string, story model.Story) error {
	sql := "UPDATE stories SET title = $1, description = $2 WHERE id = $3 AND user_id = $4"

	_, err := database.DB.Exec(sql, story.Title, story.Description, id, story.UserId)

	if err != nil {
		return err
	}
	return nil
}

func RemoveStory(id string, userId string) error {
	sql := "DELETE FROM stories WHERE id = $1 AND user_id = $2"

	_, err := database.DB.Exec(sql, id, userId)

	if err != nil {
		return err
	}
	return nil
}

func GetStoriesTask(storyId string) ([]model.Task, error) {
	sql := `SELECT
				t.*,
				COALESCE(
					SUM(
						EXTRACT(EPOCH FROM (te.end_time - te.start_time))
					),
					0
				)
				+
				CASE
					WHEN t.is_timer_running = true
						AND t.current_timer_start IS NOT NULL
					THEN EXTRACT(EPOCH FROM (NOW() - t.current_timer_start))
					ELSE 0
				END
				AS spent_hours
			FROM tasks t
			LEFT JOIN time_entries te
				ON te.task_id = t.id
			WHERE t.story_id = $1
			GROUP BY t.id;`

	rows, err := database.DB.Query(sql, storyId)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []model.Task{}

	for rows.Next() {
		var task model.Task
		err := rows.Scan(
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
			&task.SpentHours,
		)

		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
