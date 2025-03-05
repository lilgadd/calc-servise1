package router

import (
	"calc-service/internal/orchestrator/handlers"
	"github.com/gorilla/mux"
	"net/http"
)


func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Привязываем хендлеры
	r.HandleFunc("/api/v1/calculate", handlers.AddExpression).Methods("POST")
	r.HandleFunc("/internal/task", handlers.GetTaskToAgent).Methods("GET")
	r.HandleFunc("/api/v1/expressions", handlers.GetAllTasks).Methods("GET")
	r.HandleFunc("/api/v1/expressions/:id", handlers.GetExpressionByID).Methods("GET")
	r.HandleFunc("/internal/task", handlers.HandleTask).Methods("POST")
	r.HandleFunc("/api/v1/calculate", handlers.HandleCalculate).Methods("POST")
	return r
}

// StartServer запускает сервер с роутером
func StartServer() {
	r := SetupRoutes()
	http.ListenAndServe(":8080", r)
}
