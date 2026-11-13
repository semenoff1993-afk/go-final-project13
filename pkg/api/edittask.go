package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
)

// getTaskHandler возвращает задачу по ID (для формы редактирования)
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задачи"})
		}
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// updateTaskHandler сохраняет изменения задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// 1. Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	// 2. Проверяем обязательные поля
	if task.ID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// 3. Обрабатываем дату через общую функцию
	processedDate, err := processDate(task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	task.Date = processedDate

	// 4. Обновляем задачу в БД
	err = db.UpdateTask(&task)
	if err != nil {
		if errors.Is(err, db.ErrTaskNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении задачи"})
		}
		return
	}

	// 5. Возвращаем успешный ответ (пустой JSON)
	writeJSON(w, http.StatusOK, map[string]string{})
}
