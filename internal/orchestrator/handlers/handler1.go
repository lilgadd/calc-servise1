package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Структура для получения выражения из запроса
type RequestPayload struct {
	Expression string `json:"expression"`
}

// Функция для обработки POST запроса на /api/v1/calculate
func HandleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var payload RequestPayload

		// Декодирование данных из запроса
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
			return
		}

		// Выводим полученное выражение
		fmt.Printf("Получено выражение: %s\n", payload.Expression)

		// Здесь можно добавить логику для отправки задания агенту

		// Отправляем успешный ответ с уникальным ID (например, временный)
		w.WriteHeader(http.StatusCreated)
		response := map[string]string{
			"id": "1234", // Здесь будет ID задания, если ты его генерируешь
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
