package main

import (
	"fmt"
	"calc-service/router"
	"calc-service/internal/orchestrator/handlers"
	"net/http"
)

func main() {
	// Запускаем оркестратор в отдельной горутине
	go handlers.StartServer()

	// Запускаем роутер с хендлерами
	r := router.SetupRoutes()

	// Запуск сервера с роутером на порту 8080
	fmt.Println("Сервер с роутером запущен на порту :8080")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера с роутером:", err)
	}
}