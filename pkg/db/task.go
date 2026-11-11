package db

import (
	"fmt"
	"time"
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

// Tasks возвращает список ближайших задач.
// Если передана строка search, фильтрует по дате или тексту (title/comment).
func Tasks(search string, limit int) ([]Task, error) {
	// Инициализируем пустой слайс (не nil), чтобы при отсутствии задач возвращалось [] а не null
	tasks := make([]Task, 0)

	if DB == nil {
		return tasks, fmt.Errorf("база данных не инициализирована")
	}

	var query string
	var args []any

	if search != "" {
		// Проверяем, является ли поисковый запрос датой в формате DD.MM.YYYY
		parsedDate, err := time.Parse("02.01.2006", search)
		if err == nil {
			// Это дата! Ищем точное совпадение (предварительно конвертируем в 20060102)
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date ASC LIMIT ?`
			args = []any{parsedDate.Format("20060102"), limit}
		} else {
			// Это обычный текст. Ищем вхождение в title или comment (используем % для LIKE)
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date ASC LIMIT ?`
			searchPattern := "%" + search + "%"
			args = []any{searchPattern, searchPattern, limit}
		}
	} else {
		// Обычный запрос без поиска
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT ?`
		args = []any{limit}
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		var idInt int64
		// Сканируем поля. ID в базе INTEGER, но в структуре у нас string, поэтому используем временную переменную.
		err := rows.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return tasks, err
		}
		t.ID = fmt.Sprintf("%d", idInt)
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return tasks, err
	}

	return tasks, nil
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (Task, error) {
	var t Task
	if DB == nil {
		return t, fmt.Errorf("база данных не инициализирована")
	}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := DB.QueryRow(query, id)

	var idInt int64
	err := row.Scan(&idInt, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		return t, err // Если запись не найдена, вернётся sql.ErrNoRows
	}
	t.ID = fmt.Sprintf("%d", idInt)

	return t, nil
}

// UpdateTask обновляет существующую задачу в БД
func UpdateTask(task *Task) error {
	if DB == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
