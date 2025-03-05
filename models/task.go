package models

// Структура задачи
type Task struct {
    ID            string  `json:"id"`            // Уникальный идентификатор задачи
    Arg1          float64  `json:"arg1"`          // Первый аргумент для операции
    Arg2          float64  `json:"arg2"`          // Второй аргумент для операции
    Operation     string  `json:"operation"`     // Операция (например, "+", "-", "*", "/")
    OperationTime int     `json:"operation_time"` // Время, необходимое для выполнения операции
    Status string `json:"status"`                  // Статус задачи
    Result float64 `json:"result"`                   // Результат вычислений задачи
    Node          *ASTNode `json:"node,omitempty"` // Ссылка на узел дерева, с которым связана задача
}
