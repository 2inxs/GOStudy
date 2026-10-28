package main

import (
	"errors"
	"fmt"
)

type Operation func(a, b float64) (float64, error)

var operations = map[string]Operation{
	"+": func(a, b float64) (float64, error) {
		return a + b, nil
	},
	"-": func(a, b float64) (float64, error) {
		return a - b, nil
	},
	"*": func(a, b float64) (float64, error) {
		return a * b, nil
	},
	"/": func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, errors.New("деление на ноль")
		}
		return a / b, nil
	},
}

func Calculate(firstNumber float64, operator string, lastNumber float64) (float64, error) {
	operation, exists := operations[operator]

	if !exists {
		return 0, fmt.Errorf("неизвестная операция %q", operator)
	}

	return operation(firstNumber, lastNumber)
}

// я в принципе понимаю, для чего дан такой пример, но такую задачку проще делать через switch
// та же сложность O(1), по памяти одинаково
// Кода меньше и решение проще. Ну это моё мнение
func main() {
	result, err := Calculate(10, "+", 5)
	fmt.Printf("10 + 5 = %.2f\n", result)

	result, err = Calculate(10, "/", 4)
	fmt.Printf("10 / 4 = %.2f\n", result)

	result, err = Calculate(7, "*", 6)
	fmt.Printf("7 * 6 = %.2f\n", result)

	result, err = Calculate(1, "/", 0)
	fmt.Printf("1 / 0: %v\n", err)

	result, err = Calculate(2, "%", 3)
	fmt.Printf("2 %% 3: %v\n", err)
}
