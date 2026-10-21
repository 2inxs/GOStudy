package main

import "fmt"

const (
	YearInMonth = 12
	Rate        = 0.12
)

/*
	про вопрос. Деньги принято держать в копейках или центах.

float64 все равно не точно представляет большие числа после запятой
поэтому принято сначала приводить к наименьшему, а затем для вывода приводить к нужному виду
*/
func bankIncome(deposit float64) {
	var tick int
	balance := deposit
	monthlyRate := Rate / YearInMonth

	for i := range YearInMonth {
		income := balance * monthlyRate
		balance += income

		fmt.Printf("месяц %d: %.2f\n", i+1, balance)
	}

	balance = deposit
	targetBalance := deposit * 2

	for balance < targetBalance {
		income := balance * monthlyRate
		balance += income
		tick++
	}

	fmt.Printf(
		"Вклад удвоится через %d мес. Баланс: %.2f\n",
		tick,
		balance,
	)
}

func main() {
	bankIncome(30000)
}
