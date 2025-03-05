package models

// Expression представляет арифметическое выражение
type Expression struct {
    ID         string `json:"id"`         // Уникальный идентификатор выражения
    Expression string `json:"expression"` // Само арифметическое выражение
    Status     string `json:"status"`     // Статус вычисления выражения (например, "обрабатывается", "завершено")
    Result     float64 `json:"result,omitempty"` // Результат вычисления
    AST    *ASTNode `json:"-"`  // Поле для хранения дерева AST
}

// ExpressionStatus используется для возврата статуса и результата вычисления выражения
type ExpressionStatus struct {
    ID     string `json:"id"`     // Уникальный идентификатор выражения
    Status string `json:"status"` // Статус вычисления выражения
    Result float64 `json:"result,omitempty"` // Результат вычисления (опционально)
}

// Структура для ответа с Id выражения
type CreateExpressionResponse struct{
    ID string `json:"id"`
}
// Структура для ответа со списком выражений
type ListExpressionsResponse struct {
    Expressions []ExpressionStatus `json:"expressions"` // Список выражений с их статусами и результатами
}

// ExpressionTree — узел дерева выражений
type ExpressionTree struct {
    Value    float64        // Число, если это лист
    Operator string         // Оператор, если это узел операции
    Left     *ExpressionTree // Левый потомок (может быть nil)
    Right    *ExpressionTree // Правый потомок (может быть nil)
}

type TaskResult struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}