/*
К ТЗ подходит.
Убрал тогда перевод из функции,
использую перевод из копеек в рубли только при выводе
*/
package main

import "fmt"

const kopPerRub = 100

func NewAccount(initial int64) (deposit func(int64), withdraw func(int64) error, balance func() int64) {
	moneyValue := initial

	deposit = func(amount int64) {
		moneyValue = moneyValue + amount
	}

	withdraw = func(amount int64) error {
		if moneyValue < amount {
			return fmt.Errorf("недостаточно средств")
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
	depositA, withdrawA, balanceA := NewAccount(100000)
	_, _, balanceB := NewAccount(0)

	fmt.Printf("Счёт A: %.2f\n", float64(balanceA())/kopPerRub)

	depositA(50000)
	fmt.Printf("Пополнение 500.00 -> %.2f\n", float64(balanceA())/kopPerRub)

	err := withdrawA(200000)
	if err != nil {
		fmt.Printf("Снятие 2000.00: %v (баланс %.2f)\n",
			err, float64(balanceA())/kopPerRub)
	}

	err = withdrawA(30050)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Снятие 300.50 -> %.2f\n", float64(balanceA())/kopPerRub)
	}

	fmt.Printf("Счёт B: %.2f\n", float64(balanceB())/kopPerRub)
}
