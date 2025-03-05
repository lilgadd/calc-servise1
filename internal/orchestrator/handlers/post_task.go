package handlers

import (
    "encoding/json"
    "net/http"
    "calc-service/internal/orchestrator/storage"
)

// Обработчик для получения задачи
func GetTaskToAgent(w http.ResponseWriter, r *http.Request) {
    // Получаем задачу из очереди задач
    if len(storage.TaskQueue) == 0 {
        http.Error(w, "очередь задач пуста", http.StatusNotFound)
        return
    }

    // Берем первую задачу из очереди
    task := storage.TaskQueue[0]

    // Отправляем задачу в ответе
    response := map[string]interface{}{
        "task": map[string]interface{}{
            "id":            task.ID,
            "arg1":          task.Arg1,
            "arg2":          task.Arg2,
            "operation":     task.Operation,
            "operation_time": task.OperationTime,
        },
    }

    // Устанавливаем заголовки ответа и отправляем
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)

    // Удаляем задачу из очереди после того, как она была отправлена
    storage.TaskQueue = storage.TaskQueue[1:]
}
