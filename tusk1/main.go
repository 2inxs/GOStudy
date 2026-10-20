package main

import (
	"fmt"
)

const N = 5

func main() {
	buildBoard(N)
	deliveryValue(5000, 2)
	bankIncome(30000)
	uptime(5827328)
	weather(30)
}
func buildBoard(N int) {
	for i := range N {
		for j := range N {
			if (i+j)%2 == 0 {
				fmt.Print(" ")
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

const (
	Regular = iota
	Silver
	Gold
)
const DeliveryPrice = 299
const FreeDelivery = 3000

func deliveryValue(price float64, CustomerStatus int) {
	var discount float64
	var deliveryPrice float64
	var statusName string

	switch CustomerStatus {
	case Regular:
		discount = 0
		statusName = "Regular"
	case Silver:
		discount = 0.05
		statusName = "Silver"
	case Gold:
		discount = 0.1
		statusName = "Gold"
	default:
		fmt.Println("неизвестный статус")
		discount = 0
		statusName = "неизвестный статус"
	}

	discountAmount := price * discount
	priceAfterDiscount := price - discountAmount

	if priceAfterDiscount >= FreeDelivery {
		deliveryPrice = 0
	} else {
		deliveryPrice = DeliveryPrice
	}

	total := priceAfterDiscount + deliveryPrice

	fmt.Printf("Сумма: %.2f, статус: %s\n", price, statusName)
	fmt.Printf("Скидка: %.2f\n", discountAmount)
	fmt.Printf("Доставка: %.2f\n", deliveryPrice)
	fmt.Printf("Итого: %.2f\n", total)
}

const (
	YearMonth = 12
	Rate      = 0.12
)

/*
	про вопрос. Деньги принято держать в копейках или центах.

float64 все равно не точно представляет большие числа после запятой
поэтому принято сначала приводить к наименьшему, а затем для вывода приводить к нужному виду
*/
func bankIncome(fee float64) {
	var tick int
	balance := fee
	monthlyRate := Rate / YearMonth

	for i := range YearMonth {
		income := balance * monthlyRate
		balance += income

		fmt.Printf("месяц %d: %.2f\n", i+1, balance)
	}

	balance = fee
	targetBalance := fee * 2

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

func uptime(seconds int) {
	if seconds == 0 {
		fmt.Println("0с")
		return
	}
	days := seconds / (24 * 60 * 60)
	remaining := seconds % (24 * 60 * 60)

	hours := remaining / 3600
	remaining = remaining % 3600

	minutes := remaining / 60
	seconds = remaining % 60

	if days > 0 {
		fmt.Printf("%d d", days)
	}
	if hours > 0 {
		fmt.Printf("%d h", hours)
	}
	if minutes > 0 {
		fmt.Printf("%d m", minutes)
	}
	if seconds > 0 {
		fmt.Printf("%d s", seconds)
	}
}

func weather(celsius int) {
	for celsius := -20; celsius <= 40; celsius += 10 {
		fahrenheit := float64(celsius)*9/5 + 32
		var message string

		switch {
		case celsius < 0:
			message = "мороз"
		case celsius < 15:
			message = "холодно"
		case celsius < 25:
			message = "тепло"
		default:
			message = "жарко"
		}

		fmt.Printf("%d°C = %.1f°F  %s\n",
			celsius, fahrenheit, message)
	}
}
