package main

import (
	"log"
	"os"

	"github.com/semenoff1993-afk/go-final-project13/pkg/db"
	"github.com/semenoff1993-afk/go-final-project13/pkg/server"
)

func main() {
	// Задание со звёздочкой: получаем путь к БД из переменной окружения
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db" // Путь по умолчанию
	}

	// Инициализируем базу данных
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	// Запускаем веб-сервер
	server.Start()
}