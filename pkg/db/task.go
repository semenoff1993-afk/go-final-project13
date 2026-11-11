package db

import (
	"fmt"
)

// Task описывает структуру задачи для работы с JSON и БД
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в БД и возвращает её ID
func AddTask(task *Task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}