package main

import (
	"fmt"
)

const N = 5

func buildBoard() {
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

func main() {
	buildBoard()
}
