package server

import (
	"log"
	"net/http"
	"os"

	"github.com/semenoff1993-afk/go-final-project13/pkg/api" // <- Добавили импорт API
)

// Start настраивает и запускает веб-сервер
func Start() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// Регистрируем API-обработчики
	api.Init() // <- Добавили вызов Init из API

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}