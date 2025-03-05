package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"calc-service/internal/orchestrator/handlers"
)

func TestAddExpression(t *testing.T) {
	tests := []struct {
		name           string
		requestBody   string
		expectedStatus int
	}{
		{
			name:           "Valid expression",
			requestBody:    `{"expression": "3+5"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"expression": 3+5}`,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid expression",
			requestBody:    `{"expression": "3++5"}`,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte(tt.requestBody)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Создаем тестовый HTTP Recorder
			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.AddExpression)

			// Выполняем запрос
			handler.ServeHTTP(rec, req)

			// Проверяем статус код
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
