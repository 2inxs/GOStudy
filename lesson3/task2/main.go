package main

import (
	"fmt"
	"strings"
)

func AreAnagrams(firstText, secondText string) bool {
	firstText = strings.ToLower(strings.ReplaceAll(firstText, " ", ""))
	secondText = strings.ToLower(strings.ReplaceAll(secondText, " ", ""))

	charCount := make(map[rune]int)

	for _, char := range firstText {
		charCount[char]++
	}

	for _, char := range secondText {
		charCount[char]--
	}

	for _, count := range charCount {
		if count != 0 {
			return false
		}
	}

	return true
}

func main() {
	fmt.Println(AreAnagrams("listen", "silent"))
	fmt.Println(AreAnagrams("Eleven plus two", "Twelve plus one"))
	fmt.Println(AreAnagrams("Апельсин", "спаниель"))
	fmt.Println(AreAnagrams("Hello", "World"))
}
