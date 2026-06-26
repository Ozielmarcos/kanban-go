package repository

import (
	"errors"

	"github.com/Ozielmarcos/mytodolist/app/internal/model"
	"github.com/Ozielmarcos/mytodolist/app/pkg/database"
)

func GetTimeEntriesByStory(storyId string) ([]model.TimeEntry, error) {
	sql := `SELECT
				te.id,
				te.task_id,
				t.title,
				te.start_time,
				te.end_time,
				te.day,
				EXTRACT(EPOCH FROM (te.end_time - te.start_time)) AS duration
			FROM tasks t
			INNER JOIN time_entries te
				ON te.task_id = t.id
			WHERE t.story_id = $1
			ORDER BY te.day DESC, te.start_time DESC`

	rows, err := database.DB.Query(sql, storyId)
	if err != nil {
		return nil, err
	}

	var entries []model.TimeEntry

	for rows.Next() {
		var entry model.TimeEntry

		err := rows.Scan(
			&entry.ID,
			&entry.TaskId,
			&entry.TaskTitle,
			&entry.StartTime,
			&entry.EndTime,
			&entry.Day,
			&entry.Duration,
		)

		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func UpdateEntry(id string, entry model.TimeEntry) error {
	sql := `UPDATE 
				time_entries 
			SET 
				start_time = $1, 
				end_time = $2, 
				day = $3 
			WHERE id = $4
				AND task_id = $5`

	result, err := database.DB.Exec(
		sql,
		entry.StartTime,
		entry.EndTime,
		entry.Day,
		id,
		entry.TaskId,
	)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("apontamento não encontrado!")
	}

	return nil
}

func RemoveEntry(id string) error {
	sql := `DELETE FROM time_entries WHERE id = $1`

	_, err := database.DB.Exec(sql, id)
	if err != nil {
		return err
	}

	return nil
}
