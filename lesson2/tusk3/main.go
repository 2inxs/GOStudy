package main

import (
	"fmt"
	"time"
)

/*
ответ на вопрос: Потому что 2 раза переносим n-1 дисков и один раз самый большой
*/

func Hanoi(n int, from, to, via string, move func(from, to string)) {
	if n == 0 {
		return
	}

	Hanoi(n-1, from, via, to, move)
	move(from, to)
	Hanoi(n-1, via, to, from, move)
}

func RunHanoi(n int) {
	count := 0

	if n >= 20 {
		start := time.Now()

		defer func() {
			fmt.Printf("время: %v\n", time.Since(start))
		}()
	}

	move := func(from, to string) {
		count++

		if n < 20 {
			fmt.Printf("%d: %s -> %s\n", count, from, to)
		}
	}

	Hanoi(n, "A", "C", "B", move)

	if n >= 20 {
		fmt.Printf("n=%d: %d ходов\n", n, count)
	}
}

func main() {
	RunHanoi(3)
	fmt.Println()
	RunHanoi(20)
}
