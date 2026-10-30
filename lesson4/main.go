package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO: Создайте инвентарь и добавьте:
	// TODO: - Оружие: "Меч" (урон 10, прочность 5)
	// TODO: - Броню: "Щит" (защита 5, вес 4.5)
	// TODO: - Зелье: "Лечебное" (+50 HP, 3 заряда)
	// TODO: - Оружие: "Сломанный лук" (урон 5, прочность 0)

	sword := &Weapon{
		name: "Меч", level: 2, weight: 3.0, Durability: 5, quality: Common,
	}
	damage, err := sword.CalculateDamage()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка расчёта урона:", err)
		return
	}
	sword.Damage = damage

	armor := &Armor{
		name: "Кольчуга", level: 2, weight: 12.6, quality: Rare,
	}
	defense, err := armor.CalculateDefense()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка расчёта защиты:", err)
		return
	}
	armor.Defense = defense

	shield := &Armor{name: "Щит", weight: 4.5, Defense: 5}
	potion := &Potion{name: "Лечебное", weight: 0.5, Effect: "+50 HP", Charges: 3}
	bow := &Weapon{name: "Сломанный лук", weight: 2, Damage: 5, Durability: 0}

	inventory := &Inventory{}
	inventory.AddItem(sword)
	inventory.AddItem(armor)
	inventory.AddItem(shield)
	inventory.AddItem(potion)
	inventory.AddItem(bow)

	// TODO: Реализуйте логику/вызовы:
	// TODO: 1. Use предмета с выводом в консоль
	for _, item := range []Usable{sword, armor, shield, potion, bow} {
		message, err := item.Use()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Ошибка использования:", err)
			continue
		}
		fmt.Println(message)
	}

	// TODO: 2. DescribeItem с предметом, с nil и с nil-указателем (var weapon *Weapon), так же с выводом в консоль
	var weapon *Weapon
	fmt.Println(DescribeItem(sword))
	fmt.Println(DescribeItem(nil))
	fmt.Println(DescribeItem(weapon))

	// TODO: 3. Вывести в консоль результат вызова Weapons (должны вернуться только меч и лук)
	fmt.Println("Оружие:", Map(inventory.Weapons(), func(w *Weapon) string { return w.Name() }))
	// TODO: 4. Вывести в консоль результат вызова BrokenItems (должен вернуть сломанный лук)
	fmt.Println("Сломанные и израсходованные:", Map(inventory.BrokenItems(), DescribeItem))
	// TODO: 5. Вывести в консоль результат вызова ItemNames (все названия)
	fmt.Println("Названия:", inventory.ItemNames())
	// TODO: 6. Вывести в консоль результат вызова FindItemByName (поиск "Щит")
	item, found := inventory.FindItemByName("Щит")
	fmt.Println("Поиск щита:", DescribeItem(item), "найден:", found)

	// TODO: Бонус: сделайте сохранение инвентаря в файл и загрузку инвентаря из файла
	file, err := os.CreateTemp("", "gostudy-inventory-*.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка создания файла:", err)
		return
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if err := inventory.Save(file); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка сохранения инвентаря:", err)
		return
	}
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка перемещения в начало файла:", err)
		return
	}
	loaded := &Inventory{}
	if err := loaded.Load(file); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка загрузки инвентаря:", err)
		return
	}
	fmt.Println("Загружено из файла:", loaded.ItemNames())
}
