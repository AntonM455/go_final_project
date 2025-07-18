package db

import (
	"errors"
	"fmt"
)

// Task structure
type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask adds a task to the scheduler table
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, errors.New("database not initialized")
	}

	query := `
	INSERT INTO scheduler (date, title, comment, repeat)
	VALUES (?, ?, ?, ?)
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Tasks retrieves up to `limit` tasks
// from the database ordered by date ascending.
func Tasks(limit int) ([]*Task, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	query := `
		SELECT id, date, title, comment, repeat
		FROM scheduler
		ORDER BY date ASC
		LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var tasks []*Task

	for rows.Next() {
		var t Task
		var id int
		var date int
		if err := rows.Scan(&id, &date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		t.ID = fmt.Sprintf("%d", id)
		t.Date = fmt.Sprintf("%08d", date)
		tasks = append(tasks, &t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

// GetTask returns one task by id
func GetTask(id string) (*Task, error) {
	if DB == nil {
		return nil, errors.New("database not initialized")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	var t Task
	var idNum int
	var date int

	err := DB.QueryRow(query, id).Scan(&idNum, &date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return nil, err
	}

	t.ID = fmt.Sprintf("%d", idNum)
	t.Date = fmt.Sprintf("%08d", date)
	return &t, nil
}

// UpdateTask updates task fields by id
func UpdateTask(task *Task) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `
	UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ?
	WHERE id = ?
	`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// DeleteTask deletes a task by ID
func DeleteTask(id string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	// check how many rows were deleted
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// UpdateDate updates only the task due date
func UpdateDate(next string, id string) error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	// check if the row has been updated
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
