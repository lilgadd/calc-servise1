package storage

import (
	"sync"
	"calc-service/models"
	"errors"
)

var (
	TasksMap       = make(map[string]*models.Task)       // Мапа для хранения задач
	TaskQueue = make([]*models.Task, 0) 				 // Мапа для хранения результатов задач
	mu             sync.RWMutex                          // Мьютекс для синхронизации доступа
)

// Функция для получения выражения по ID с нужным форматом
func GetExpressionByID(taskID string) (map[string]models.ExpressionStatus, error) {
	// Захватываем блокировку на чтение
	mu.RLock()
	defer mu.RUnlock()

	// Проверяем, существует ли задача с данным ID
	task, exists := TasksMap[taskID]
	if !exists {
		return nil, errors.New("задача не найдена")
	}

	// Создаём нужный результат
	expressionResult := map[string]models.ExpressionStatus{
		"expression": {
			ID:     taskID,
			Status: task.Status,
			Result: task.Result,
		},
	}

	// Возвращаем результат
	return expressionResult, nil
}


// Добавление задачи в хранилище
func AddTaskToStorage(task *models.Task) {
	mu.Lock()
	defer mu.Unlock()
	TasksMap[task.ID] = task
}

// Получение задачи по ID
func GetTask(taskID string) (*models.Task, bool) {
	mu.RLock()
	defer mu.RUnlock()
	task, exists := TasksMap[taskID]
	return task, exists
}

// Удаление задачи по ID
func RemoveTask(taskID string) {
    delete(TasksMap, taskID)          // Удаляем задачу из карты
    for i, task := range TaskQueue { // Удаляем задачу из очереди
        if task.ID == taskID {
            TaskQueue = append(TaskQueue[:i], TaskQueue[i+1:]...)
            break
        }
    }
}
