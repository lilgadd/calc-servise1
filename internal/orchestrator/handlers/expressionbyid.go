package handlers

import (
	"calc-service/internal/orchestrator/storage"
	"calc-service/models"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

// Получение выражения по идентификатору
func GetExpressionByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Получаем выражение из хранилища по ID
	expression, exists := storage.GetTask(id)
	if !exists {
		// Если выражения с таким ID нет, возвращаем 404
		http.Error(w, "выражение не найдено", http.StatusNotFound)
		return
	}

	// Формируем ответ
	response := models.ExpressionStatus{
		ID: expression.ID,
		Status: expression.Status,
		Result: expression.Result,
	}

	// Обработка ошибок, если что-то пошло не так при отправке ответа
	defer func() {
		if err := recover(); err != nil {
			// Отправляем статус 500, если произошла ошибка
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}()
	
	// Устанавливаем заголовок и код ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 - успешно
	json.NewEncoder(w).Encode(response)
}
