// Ответ на вопрос: append не создает новый слайс, а работает в исходном
// на который ссылается nums
// поэтому сначала создаем новый, а потом копируем
package main

import "fmt"

func RemoveElement(nums []int, index int) ([]int, error) {
	if index < 0 || index >= len(nums) {
		return nil, fmt.Errorf("индекс %d вне диапазона [0, %d)", index, len(nums))
	}

	result := make([]int, 0, len(nums)-1)

	result = append(result, nums[:index]...)
	result = append(result, nums[index+1:]...)

	return result, nil
}

func main() {
	nums := []int{1, 2, 2, 3, 1}

	result, err := RemoveElement(nums, 3)
	fmt.Println(result, err)
	fmt.Println("Исходный слайс:", nums)

	result, err = RemoveElement([]int{1, 2, 3}, 5)
	fmt.Println(result, err)

	result, err = RemoveElement([]int{1, 2, 3}, -1)
	fmt.Println(result, err)
}
