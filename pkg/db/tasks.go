package db

import "database/sql"

func Tasks(limit int) ([]*Task, error) {

	var tasks []*Task

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit", sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		task := &Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return []*Task{}, nil
	}
	return tasks, nil

}
