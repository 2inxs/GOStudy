package main

import (
	"fmt"
	"time"
)

/*
Ответ на вопрос: потому что дважды переносим n-1 дисков
и один раз самый большой.
*/

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}

	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
}

func MeasureTime() func() {
	start := time.Now()

	return func() {
		fmt.Printf("время: %v\n", time.Since(start))
	}
}

func main() {
	// n = 3
	count := 0

	printMove := func(from, to string) {
		count++
		fmt.Printf("%d: %s -> %s\n", count, from, to)
	}

	Hanoi(3, "A", "C", "B", printMove)

	fmt.Println()

	// n = 20
	count = 0

	countMove := func(from, to string) {
		count++
	}

	defer MeasureTime()()

	Hanoi(20, "A", "C", "B", countMove)

	fmt.Printf("n=20: %d ходов\n", count)
}
