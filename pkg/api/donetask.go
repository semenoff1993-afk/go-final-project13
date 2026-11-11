package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
)

// doneTaskHandler обрабатывает POST-запрос /api/task/done
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// 1. Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	// 2. Если правило пустое — удаляем задачу
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении задачи"})
			return
		}
	} else {
		// 3. Если периодическая — вычисляем следующую дату
		next, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка вычисления даты: %v", err)})
			return
		}

		// Обновляем дату и сохраняем обратно в БД
		task.Date = next
		err = db.UpdateTask(&task)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при обновлении задачи"})
			return
		}
	}

	// Возвращаем пустой JSON при успехе
	writeJSON(w, http.StatusOK, map[string]string{})
}

// deleteTaskHandler обрабатывает DELETE-запрос /api/task
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		if err.Error() == "задача не найдена" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при удалении задачи"})
		}
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
