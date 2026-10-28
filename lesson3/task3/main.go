// Мои любимые задачки с leetcode)))
package main

import "fmt"

func FirstUnique(text string) rune {
	charCount := make(map[rune]int)

	for _, char := range text {
		charCount[char]++
	}

	for _, char := range text {
		if charCount[char] == 1 {
			return char
		}
	}

	return 0
}

func main() {
	fmt.Printf("%q\n", FirstUnique("abacabad"))
	fmt.Printf("%q\n", FirstUnique("abcabcd"))
	fmt.Printf("%q\n", FirstUnique("абракадабра"))
	fmt.Printf("%q\n", FirstUnique("aabbcc"))
}
