/*
Не совсем понял логику ТЗ.
По условию деньги хранятся в копейках (int64), но выводить их нужно в рублях с двумя знаками после точки.
При этом balance возвращает int64.

На мой взгляд, логичнее оставить всю бизнес-логику на int64, включая входные параметры и возвращаемые значения.
В реальном backend нет смысла преобразовывать суммы в float64 ради вывода через fmt.Printf.
Проще передавать на frontend копейки, а форматировать их в рубли уже на стороне клиента.

Если же форматирование должно происходить внутри программы, логичнее вынести его в отдельную функцию, не смешивая с получением баланса.
*/
package main

import "fmt"

const kopPerRub = 100

func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	moneyValue := initial * kopPerRub

	deposit = func(amount int64) {
		moneyValue = moneyValue + amount*kopPerRub
	}

	withdraw = func(amount int64) error {
		amount *= kopPerRub

		if moneyValue < amount {
			return fmt.Errorf(
				"недостаточно средств (баланс %.2f)",
				float64(moneyValue)/kopPerRub,
			)
		}

		moneyValue = moneyValue - amount
		return nil
	}

	balance = func() int64 {
		return moneyValue
	}

	return deposit, withdraw, balance
}

func main() {
	depositA, withdrawA, balanceA := NewAccount(1000)
	_, _, balanceB := NewAccount(0)

	fmt.Printf("Счёт A: %.2f\n", float64(balanceA())/kopPerRub)

	depositA(500)
	fmt.Printf("Пополнение 500.00 -> %.2f\n", float64(balanceA())/kopPerRub)

	err := withdrawA(2000)
	if err != nil {
		fmt.Printf("Снятие 2000.00: %v\n", err)
	}

	err = withdrawA(300)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Снятие 300.00 -> %.2f\n", float64(balanceA())/kopPerRub)
	}

	fmt.Printf("Счёт B: %.2f\n", float64(balanceB())/kopPerRub)
}
