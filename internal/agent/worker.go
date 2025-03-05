package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"calc-service/internal/orchestrator/storage"
)

// Структура для отправки результата в оркестратор
type ExpressionPayload struct {
	Expression struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	} `json:"expression"`
}

// Функция для отправки результата в оркестратор
func sendResult(taskID string, result float64) error {
	// Используем структуру ExpressionPayload
	payload := ExpressionPayload{
		Expression: struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}{
			ID:     taskID,
			Result: result,
		},
	}

	// Сериализация данных в JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных результата: %v", err)
	}

	// Отправка результата оркестратору
	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("ошибка в отправке результата: %v", err)
	}
	defer resp.Body.Close()

	// Обработка ответа от оркестратора
	switch resp.StatusCode {
	case http.StatusOK:
		log.Printf("Результат задачи %s успешно сохранен", taskID)
		return nil
	case http.StatusNotFound:
		log.Printf("Задача с ID: %s не найдена", taskID)
		return fmt.Errorf("задача не найдена")
	case http.StatusUnprocessableEntity:
		log.Printf("Некорректный ID задачи %s", taskID)
		return fmt.Errorf("некорректные данные")
	case http.StatusInternalServerError:
		log.Printf("Ошибка сервера во время обработки: %s", taskID)
		return fmt.Errorf("ошибка сервера")
	default:
		log.Printf("Неизвестный код ответа %d для задачи %s", resp.StatusCode, taskID)
		return fmt.Errorf("неизвестный код ответа %d", resp.StatusCode)
	}
}

// Функция воркера для обработки задачи
func worker(id int) {
	for {
		// Получаем ID задачи
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Printf("Worker %d: ошибка получения задачи: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Если задачи нет, ждем и пробуем снова
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		var taskResp struct {
			ID string `json:"id"`
		}

		// Декодируем ID задачи из ответа
		err = json.NewDecoder(resp.Body).Decode(&taskResp)
		resp.Body.Close() // Закрываем тело ответа

		if err != nil {
			log.Printf("Worker %d: ошибка декодирования ID задачи: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Обрабатываем задачу
		taskID := taskResp.ID
		storage.TasksMap[taskID].Status = "выполняется"
		result, err := calculate(taskID)
		storage.TasksMap[taskID].Status = "завершен"
		storage.TasksMap[taskID].Result = result 
		if err != nil {
			log.Printf("Worker %d: ошибка при вычислении задачи %s: %v", id, taskID, err)
			result = 0
			storage.TasksMap[taskID].Status = "ошибка"
			time.Sleep(1 * time.Second)
			continue
		}

		// Отправляем результат
		err = sendResult(taskID, result)
		if err != nil {
			log.Printf("Worker %d: ошибка отправки результата для задачи %s: %v", id, taskID, err)
			time.Sleep(1 * time.Second)
			continue
		}

		// Пауза перед следующим циклом обработки
		time.Sleep(1 * time.Second)
	}
}
