package main

import "fmt"

type status int

const (
	Regular status = iota
	Silver
	Gold
)
const (
	DeliveryPrice = 299
	FreeDelivery  = 3000
)

func deliveryValue(price float64, customerStatus status) {
	var discount float64
	var deliveryPrice float64
	var statusName string

	switch customerStatus {
	case Regular:
		statusName = "Regular"
	case Silver:
		discount = 0.05
		statusName = "Silver"
	case Gold:
		discount = 0.1
		statusName = "Gold"
	default:
		statusName = "неизвестный статус"
	}

	discountAmount := price * discount
	priceAfterDiscount := price - discountAmount

	if priceAfterDiscount >= FreeDelivery {
		deliveryPrice = 0
	} else {
		deliveryPrice = DeliveryPrice
	}

	total := priceAfterDiscount + deliveryPrice

	fmt.Printf("Сумма: %.2f, статус: %s\n", price, statusName)
	fmt.Printf("Скидка: %.2f\n", discountAmount)
	fmt.Printf("Доставка: %.2f\n", deliveryPrice)
	fmt.Printf("Итого: %.2f\n", total)
}

func main() {
	deliveryValue(5000, Gold)
}
