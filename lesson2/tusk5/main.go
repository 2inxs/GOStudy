package main

import "fmt"

func CountBits(n uint8) int {
	count := 0

	for n > 0 {
		if n&1 == 1 {
			count++
		}

		n = n >> 1
	}

	return count
}

func IsPowerOfTwo(n uint) bool {
	return n != 0 && (n&(n-1)) == 0
}

func main() {
	numbers := []uint8{0, 5, 12, 255}

	for _, n := range numbers {
		fmt.Printf("%d = %08b -> %d\n", n, n, CountBits(n))
	}

	fmt.Println()

	values := []uint{1, 6, 64, 0}

	for _, n := range values {
		fmt.Printf("%d -> %t\n", n, IsPowerOfTwo(n))
	}
}
