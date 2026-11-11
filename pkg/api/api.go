package api

import (
	"net/http"
)

// DateFormat определяет стандартный формат даты для проекта
const DateFormat = "20060102"

// Init регистрирует все HTTP-обработчики
func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
}