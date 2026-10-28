package main

import (
	"fmt"
	"strings"
)

func WordCount(s string) map[string]int {
	words := strings.Fields(s)
	wordCount := make(map[string]int)

	for _, word := range words {
		wordCount[word]++
	}

	return wordCount
}

func main() {
	fmt.Println(WordCount("go is fun go"))
	fmt.Println(WordCount("Hello hello World"))
	fmt.Println(WordCount(" мир мир труд "))
}
