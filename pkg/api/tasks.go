package api

import (
	"net/http"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
)

// tasksHandler обрабатывает GET-запрос /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	// Получаем строку поиска, если она есть
	search := r.FormValue("search")

	// Запрашиваем задачи с лимитом 50 штук
	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка при получении задач"})
		return
	}

	// Возвращаем результат
	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}