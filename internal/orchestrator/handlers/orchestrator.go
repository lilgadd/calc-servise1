package handlers

import (
	"calc-service/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"strconv"
	
	"regexp"
	"calc-service/internal/orchestrator/storage"
	"fmt"
	"unicode"
)

var (
	expressions = make(map[string]*models.Expression) // Мапа для хранения выражений
)

// Добавление нового выражения
func AddExpression(w http.ResponseWriter, r *http.Request) {
	var expr models.Expression

	// Получаем тело запроса
	err := json.NewDecoder(r.Body).Decode(&expr)
	if err != nil {
		http.Error(w, "некорректные данные", http.StatusInternalServerError)
		return
	}

	// Проверка на валидность
	if !isValidExpression(expr.Expression) {
		http.Error(w, "некорректные данные", http.StatusUnprocessableEntity)
		return
	}

	// Очистить выражение от пробелов
	expr.Expression = strings.ReplaceAll(expr.Expression, " ", "")

	expr.ID = generateID()

	expr.Status = "ожидает выполнения"

	rpn := convertToRPN(expr.Expression)

	expressionTree := createExpressionTree(rpn)

	createTasksForTree(expressionTree)

	// Возвращаем ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 - Создано
	json.NewEncoder(w).Encode(models.CreateExpressionResponse{ID: expr.ID})
}

// Функция для проверки валидности выражения
func isValidExpression(expression string) bool {
	// Проверка на баланс скобок
	openBrackets := 0
	previousChar := ""
	for i, char := range expression {
		// Проверка на баланс скобок
		if char == '(' {
			openBrackets++
		} else if char == ')' {
			openBrackets--
		}

		// Проверка на недопустимые символы
		if !strings.Contains("0123456789+-*/().", string(char)) && !unicode.IsSpace(char) {
			return false
		}

		// Проверка на несколько подряд идущих операторов
		if (char == '+' || char == '-' || char == '*' || char == '/') && 
			(previousChar == "+" || previousChar == "-" || previousChar == "*" || previousChar == "/") {
			return false
		}

		// Проверка на неправильное начало выражения
		if i == 0 && (char == '+' || char == '*' || char == '/') {
			return false
		}

		// Проверка на неправильное завершение выражения
		if i == len(expression)-1 && (char == '+' || char == '*' || char == '/') {
			return false
		}

		previousChar = string(char)
	}

	// Проверка баланса скобок
	return openBrackets == 0
}

// Генерация ID для выражения
func generateID() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// Преобразование выражения в обратную польскую нотацию
func convertToRPN(expression string) []string {
	// Используем стэк для преобразования в обратную польскую нотацию
	var output []string
	var stack []string

	// Регулярное выражение для определения чисел
	re := regexp.MustCompile(`(\d+\.?\d*|\+|\-|\*|\/|\(|\))`)

	tokens := re.FindAllString(expression, -1)

	// Алгоритм Шёнхеймера для конвертации в RPN
	for _, token := range tokens {
		switch token {
		case "+", "-":
			for len(stack) > 0 && (stack[len(stack)-1] == "+" || stack[len(stack)-1] == "-" || stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "*", "/":
			for len(stack) > 0 && (stack[len(stack)-1] == "*" || stack[len(stack)-1] == "/") {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1] // Убираем "("
		default:
			output = append(output, token)
		}
	}

	// Добавляем оставшиеся операторы в стек в конец
	for len(stack) > 0 {
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output
}

// Создание дерева выражений
func createExpressionTree(rpn []string) *models.ASTNode {
    var stack []*models.ASTNode // Стек для построения дерева

    for _, token := range rpn {
        if token == "+" || token == "-" || token == "*" || token == "/" {
            // Если оператор, извлекаем два операнда
            right := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            left := stack[len(stack)-1]
            stack = stack[:len(stack)-1]

            // Создаем новый узел с этим оператором
            node := &models.ASTNode{
                Operator: token,
                Left:     left,
                Right:    right,
            }

            // Помещаем узел обратно в стек
            stack = append(stack, node)
        } else {
            // Если число, создаем лист
            value := 0.0
            fmt.Sscanf(token, "%f", &value) // Преобразуем строку в число

            // Создаем новый узел с числом
            node := &models.ASTNode{
                Value:   value,
                IsLeaf:  true, // Лист — это число
                Operator: "",
            }

            // Помещаем узел обратно в стек
            stack = append(stack, node)
        }
    }

    // В стеке останется один элемент — это корень дерева
    return stack[0]
}

// Создание задач для обработки дерева

var taskid int

// Функция для планирования задач, основанных на дереве выражений
func createTasksForTree(node *models.ASTNode){
    var tasks []*models.Task
    var traverse func(node *models.ASTNode)
    traverse = func(node *models.ASTNode) {
        if node == nil || node.IsLeaf {
            return
        }
        traverse(node.Left)
        traverse(node.Right)
        
        if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
            if !node.TaskScheduled {
                taskid++
                taskID := fmt.Sprintf("%d", taskid)
                var opTime int
                
				opTime = getOperationTime(node.Operator)

                task := &models.Task{
                    ID:            taskID,
                    Arg1:          node.Left.Value,  
    				Arg2:          node.Right.Value,
                    Operation:     node.Operator,
                    OperationTime: opTime,
                    Node:          node,
                }
                node.TaskScheduled = true
                storage.TasksMap[taskID] = task
                storage.TaskQueue = append(storage.TaskQueue, task)

                // Добавляем задачу в список
                tasks = append(tasks, task)
            }
        }
    }
    traverse(node)
}




// Функция, которая возвращает время выполнения операции
func getOperationTime(operator string) int {
    switch operator {
    case "+":
        return models.TimeAddition 
    case "-":
        return models.TimeSubtraction 
    case "*":
        return models.TimeMultiplication
    case "/":
        return models.TimeDivision
    default:
        return 0
    }
}


// Запуск сервера
func StartServer(){
	// Запуск сервера
	fmt.Println("Сервер оркестратора запущен на порту :8080")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера оркестратора:", err)
	}
}