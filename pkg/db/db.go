package db

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// DB — глобальная переменная для подключения к базе данных
var DB *sql.DB

// Init открывает базу данных и создает таблицу, если файла еще нет
func Init(dbFile string) error {
	install := false

	// Проверяем, существует ли файл
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		install = true
	}

	// Открываем подключение
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// Если файла не было, создаем таблицу и индекс
	if install {
		schema := `
		CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT "",
			title VARCHAR(256) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);
		CREATE INDEX scheduler_date_idx ON scheduler (date);
		`
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
		log.Println("Файл базы данных и таблица scheduler успешно созданы.")
	}

	return nil
}