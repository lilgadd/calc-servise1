package handlers

import (
	"calc-service/internal/orchestrator/storage"
	"calc-service/models"
	"encoding/json"
	"net/http"
	"sync"
)

var mu sync.RWMutex

// Получение всего списка выражений
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Захватываем блокировку на чтение
	mu.RLock()
	defer mu.RUnlock()

	var expressionList []models.ExpressionStatus

	// Проходим по всем задачам в TaskMap
	for id, task := range storage.TasksMap {
		// Создаём список ExpressionStatus из данных задачи
		expressionList = append(expressionList, models.ExpressionStatus{
			ID:     id,
			Status: task.Status,
			Result: task.Result,
		})
	}

	// Если задач нет, возвращаем 204 No Content
	if len(expressionList) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем список задач в формате JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 - ОК
	if err := json.NewEncoder(w).Encode(map[string][]models.ExpressionStatus{"expressions": expressionList}); err != nil {
		http.Error(w, "Ошибка при кодировании данных", http.StatusInternalServerError)
		return
	}
}
