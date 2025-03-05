package agent

import (
	"calc-service/internal/orchestrator/storage"
	"fmt"
	"time"
)

// Функция для выполнения вычислений
func calculate(taskID string) (float64, error) {
	// Получаем задачу из TaskMap
	task, exists := storage.TasksMap[taskID]
	if !exists {
		return 0, fmt.Errorf("задача с ID: %s не найдена", taskID)
	}

	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	// Выполняем вычисления
	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 == 0 {
			return 0, fmt.Errorf("деление на 0")
		}
		result = task.Arg1 / task.Arg2
	default:
		return 0, fmt.Errorf("неизвестная операция: %v", task.Operation)
	}

	return result, nil
}