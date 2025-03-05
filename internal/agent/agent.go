package agent

import (
	"log"
	"calc-service/internal/orchestrator/storage"
)

// RunAgent запускает агентов, равное количеству задач в TaskMap.
func RunAgent() {
	// Получаем количество задач, соответствующее количеству воркеров
	computingPower := len(storage.TasksMap)
	if computingPower == 0 {
		log.Println("Нет задач для обработки.")
		return
	}

	// Запускаем нужное количество воркеров
	for i := 0; i < computingPower; i++ {
		go worker(i)
	}

	// Ожидаем завершения работы воркеров (не завершаем выполнение программы)
	select {}
}
