package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lilgadd/calc.go/pkg/rpn"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	Config *Config
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New() *Application {
	return &Application{
		Config: ConfigFromEnv(),
	}
}

// Функция запуска приложения
// тут будем чиать введенную строку и после нажатия ENTER писать результат работы программы на экране
// если пользователь ввел exit - то останаваливаем приложение
func (a *Application) Run() error {
	for {
		log.Println("input expression")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			errorResponse := ErrorResponse{
				Error: "Failed to read expression from console",
			}
			json.NewEncoder(os.Stdout).Encode(errorResponse)
			continue
		}
		text = strings.TrimSpace(text)
		if text == "exit" {
			log.Println("application was successfully closed")
			return nil
		}
		//вычисление выражения
		result, err := calculation.Calc(text)
		if err != nil {
			errorResponse := ErrorResponse{
				Error: err.Error(),
			}
			json.NewEncoder(os.Stdout).Encode(errorResponse)
		} else {
			resultResponse := struct {
				Result float64 `json:"result"`
			}{
				Result: result,
			}
			json.NewEncoder(os.Stdout).Encode(resultResponse)
		}
	}
}

type Request struct {
	Expression string `json:"expression"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost{
		errorResponse := ErrorResponse{
			Error: "Method not allowed",
		}
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	
	request := new(Request)
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
        http.Error(w, "{\"error\": \"Invalid JSON format\"}", http.StatusBadRequest)
        return
	}

	request.Expression = strings.ReplaceAll(request.Expression, " ", "")

	result, err := calculation.Calc(request.Expression)
	if err != nil {
    	// Если ошибка связана с делением на ноль
    	if strings.Contains(err.Error(), "деление на ноль") {
        	errorResponse := map[string]string{
            	"error": "Division by zero is not allowed",
        	}
        	w.WriteHeader(http.StatusInternalServerError) // Код 500
        	json.NewEncoder(w).Encode(errorResponse)
    	} else if strings.Contains(err.Error(), "недопустимый символ") || strings.Contains(err.Error(), "несбалансированные скобки") {
        	// Ошибка синтаксиса
        	errorResponse := map[string]string{
            	"error": "Expression is not valid",
        	}
        	w.WriteHeader(http.StatusUnprocessableEntity) // Код 422
        	json.NewEncoder(w).Encode(errorResponse)
    	} else {
        	// Для остальных ошибок
        	errorResponse := map[string]string{
            	"error": "Internal server error: " + err.Error(),
        	}
        	w.WriteHeader(http.StatusInternalServerError) // Код 500
        	json.NewEncoder(w).Encode(errorResponse)
    	}
	} else {
		resultResponse := map[string]string{
			"result": fmt.Sprintf("%f", result),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resultResponse)
	}
}

func (a *Application) RunServer() error {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	return http.ListenAndServe(":"+a.Config.Addr, nil)
}