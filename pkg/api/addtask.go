package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
)

// writeJSON - вспомогательная функция для отправки JSON ответов
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ошибка при кодировании JSON: %v", err)
	}
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

	// 3. Обрабатываем дату через общую функцию
	processedDate, err := processDate(task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	task.Date = processedDate

	// 4. Добавляем в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при добавлении задачи в БД"})
		return
	}

	// 5. Возвращаем успешный ответ
	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}
