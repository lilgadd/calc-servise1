package calculation_test

import (
	"testing"
	"github.com/lilgadd/calc.go/pkg/rpn"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedResult float64
		expectedError  string
	}{
		{
			name:           "Valid expression with addition",
			expression:     "2 + 3",
			expectedResult: 5,
			expectedError:  "",
		},
		{
			name:           "Valid expression with subtraction",
			expression:     "5 - 2",
			expectedResult: 3,
			expectedError:  "",
		},
		{
			name:           "Valid expression with multiplication",
			expression:     "3 * 4",
			expectedResult: 12,
			expectedError:  "",
		},
		{
			name:           "Valid expression with division",
			expression:     "10 / 2",
			expectedResult: 5,
			expectedError:  "",
		},
		{
			name:           "Division by zero",
			expression:     "10 / 0",
			expectedResult: 0,
			expectedError:  "деление на ноль",
		},
		{
			name:           "Expression with parentheses",
			expression:     "(2 + 3) * 4",
			expectedResult: 20,
			expectedError:  "",
		},
		{
			name:           "Nested parentheses",
			expression:     "2 + (3 * (4 - 1))",
			expectedResult: 11,
			expectedError:  "",
		},
		{
			name:           "Invalid expression with unbalanced parentheses",
			expression:     "(2 + 3",
			expectedResult: 0,
			expectedError:  "несбалансированные скобки",
		},
		{
			name:           "Invalid expression with invalid character",
			expression:     "2 + $",
			expectedResult: 0,
			expectedError:  "недопустимый символ: $",
		},
		{
			name:           "Empty expression",
			expression:     "",
			expectedResult: 0,
			expectedError:  "нет чисел для вычисления",
		},
		{
			name:           "Expression with multiple operators",
			expression:     "2 + 3 * 4 - 1",
			expectedResult: 13,
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculation.Calc(tt.expression)

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err.Error())
			}

			if result != tt.expectedResult {
				t.Errorf("expected result: %v, got: %v", tt.expectedResult, result)
			}
		})
	}
}
