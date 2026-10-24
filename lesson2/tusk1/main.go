package main

import "fmt"

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
			return 0, nil
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

}
