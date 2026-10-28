package main

import (
	"fmt"
	"unicode"
)

func IsPalindrome(text string) bool {
	chars := []rune(text)

	left := 0
	right := len(chars) - 1

	for left < right {
		if !unicode.IsLetter(chars[left]) {
			left++
			continue
		}

		if !unicode.IsLetter(chars[right]) {
			right--
			continue
		}

		if unicode.ToLower(chars[left]) != unicode.ToLower(chars[right]) {
			return false
		}

		left++
		right--
	}

	return true
}

func main() {
	fmt.Println(IsPalindrome("шалаш"))
	fmt.Println(IsPalindrome("А роза упала на лапу Азора"))
	fmt.Println(IsPalindrome("Was it a car or a cat I saw?"))
	fmt.Println(IsPalindrome("hello"))
}
