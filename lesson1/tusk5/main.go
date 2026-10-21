package main

import "fmt"

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

func main() {
	weather(30)
}
