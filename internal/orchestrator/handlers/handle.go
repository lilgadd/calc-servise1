package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"calc-service/models"
)

// Функция для обработки POST запроса на /internal/task
func HandleTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var exp models.TaskResult
		// Декодирование данных из запроса
		err := json.NewDecoder(r.Body).Decode(&exp)
		if err != nil {
			http.Error(w, "Ошибка при декодировании данных", http.StatusBadRequest)
			return
		}

		// Выводим полученные данные без поля status
		fmt.Printf("Получено задание: ID=%s, Результат=%.2f\n", exp.ID, exp.Result)

		// Отправляем успешный ответ
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Результат успешно получен"))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
