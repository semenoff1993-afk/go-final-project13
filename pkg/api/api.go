package api

import (
	"fmt"
	"net/http"
	"time"
)

const DateFormat = "20060102"

// Init регистрирует все HTTP-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/signin", signinHandler)

	// Оборачиваем маршруты в middleware
	http.HandleFunc("/api/task", authMiddleware(taskHandler))
	http.HandleFunc("/api/tasks", authMiddleware(tasksHandler))
	http.HandleFunc("/api/task/done", authMiddleware(doneTaskHandler))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

// processDate проверяет и вычисляет дату задачи
func processDate(date, repeat string) (string, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	if date == "" {
		date = today.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("Дата представлена в неверном формате")
	}

	if t.Before(today) {
		if repeat == "" {
			date = today.Format(DateFormat)
		} else {
			next, err := NextDate(now, date, repeat)
			if err != nil {
				return "", fmt.Errorf("Ошибка правила повторения: %v", err)
			}
			date = next
		}
	} else if repeat != "" {
		_, err := NextDate(now, date, repeat)
		if err != nil {
			return "", fmt.Errorf("Ошибка правила повторения: %v", err)
		}
	}

	return date, nil
}
