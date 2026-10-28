package main

import "fmt"

const (
	secondsPerMinute = 60
	secondsPerHour   = 60 * secondsPerMinute
	secondsPerDay    = 24 * secondsPerHour
)

func uptime(seconds int) {
	fmt.Printf("%d -> ", seconds)

	if seconds == 0 {
		fmt.Println("0с")
		return
	}
	days := seconds / secondsPerDay
	remaining := seconds % secondsPerDay

	hours := remaining / secondsPerHour
	remaining = remaining % secondsPerHour

	minutes := remaining / secondsPerMinute
	seconds = remaining % secondsPerMinute

	separator := ""

	if days > 0 {
		fmt.Printf("%s%dд", separator, days)
		separator = " "
	}
	if hours > 0 {
		fmt.Printf("%s%dч", separator, hours)
		separator = " "
	}
	if minutes > 0 {
		fmt.Printf("%s%dм", separator, minutes)
		separator = " "
	}
	if seconds > 0 {
		fmt.Printf("%s%dс", separator, seconds)
	}

	fmt.Println()
}

func main() {
	uptime(5827328)
}
