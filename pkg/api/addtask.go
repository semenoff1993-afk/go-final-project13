package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
)

// writeJSON - вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	
	// 1. Десериализуем JSON
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	// 2. Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	now := time.Now()
	// Сбрасываем время до нулей, чтобы корректно сравнивать только даты
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 3. Обрабатываем дату
	if task.Date == "" {
		task.Date = today.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Дата представлена в неверном формате"})
		return
	}

	// Проверяем прошлые даты и правила повторения
	if t.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format(DateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка правила повторения: %v", err)})
				return
			}
			task.Date = next
		}
	} else if task.Repeat != "" {
		// Даже если дата в будущем, правило повторения нужно валидировать
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка правила повторения: %v", err)})
			return
		}
	}

	// 4. Добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении задачи в БД"})
		return
	}

	// 5. Возвращаем успешный ответ
	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}