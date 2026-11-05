package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gostudy/lesson4/internal/inventory"
	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/pkg"
)

func main() {
	if err := Run(); err != nil {
		panic(err)
	}
}

func Run() error {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	loadPath := flags.String("load", "", "файл для загрузки инвентаря")
	savePath := flags.String("save", "", "файл для сохранения инвентаря")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("неизвестные аргументы: %v", flags.Args())
	}

	var inv *inventory.Inventory
	var err error
	if *loadPath != "" {
		inv, err = loadInventory(*loadPath)
	} else {
		inv, err = newDemoInventory()
	}
	if err != nil {
		return err
	}

	for _, item := range inv.Items {
		usable, ok := item.(items.Usable)
		if !ok {
			continue
		}
		message, err := usable.Use()
		if err != nil {
			if errors.Is(err, items.ErrBrokenWeapon) || errors.Is(err, items.ErrNoCharges) {
				fmt.Fprintln(os.Stderr, "Ошибка использования:", err)
				continue
			}
			return fmt.Errorf("использование предмета: %w", err)
		}
		fmt.Println(message)
	}

	var weapon *items.Weapon
	if len(inv.Items) > 0 {
		fmt.Println(items.DescribeItem(inv.Items[0]))
	}
	fmt.Println(items.DescribeItem(nil))
	fmt.Println(items.DescribeItem(weapon))
	fmt.Println("Оружие:", helpers.Map(inv.Weapons(), func(w *items.Weapon) string { return w.Name() }))
	fmt.Println("Сломанные и израсходованные:", helpers.Map(inv.BrokenItems(), items.DescribeItem))
	fmt.Println("Названия:", inv.ItemNames())
	item, found := inv.FindItemByName("Щит")
	fmt.Println("Поиск щита:", items.DescribeItem(item), "найден:", found)

	var file *os.File
	if *savePath != "" {
		file, err = os.Create(*savePath)
	} else {
		file, err = os.CreateTemp("", "gostudy-inventory-*.txt")
	}
	if err != nil {
		return fmt.Errorf("создание файла инвентаря: %w", err)
	}
	if *savePath == "" {
		defer os.Remove(file.Name())
	}
	if err := inv.Save(file); err != nil {
		file.Close()
		return fmt.Errorf("сохранение инвентаря: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("закрытие файла инвентаря: %w", err)
	}
	loaded, err := loadInventory(file.Name())
	if err != nil {
		return err
	}
	fmt.Println("Загружено из файла:", loaded.ItemNames())
	return nil
}

func newDemoInventory() (*inventory.Inventory, error) {
	sword := items.NewWeapon(101, "Меч", 2, 3.0, items.Common)
	sword.Durability = 5
	damage, err := sword.CalculateDamage()
	if err != nil {
		return nil, fmt.Errorf("расчёт урона: %w", err)
	}
	sword.Damage = damage

	armor := items.NewArmor(201, "Кольчуга", 2, 12.6, items.Rare)
	defense, err := armor.CalculateDefense()
	if err != nil {
		return nil, fmt.Errorf("расчёт защиты: %w", err)
	}
	armor.Defense = defense

	shield := items.NewArmor(202, "Щит", 1, 4.5, items.Common)
	shield.Defense = 5
	potion := items.NewPotion(301, "Лечебное", 0.5, "+50 HP", 3)
	bow := items.NewWeapon(102, "Сломанный лук", 1, 2, items.Common)
	bow.Damage = 5

	inv := &inventory.Inventory{}
	inv.AddItem(sword)
	inv.AddItem(armor)
	inv.AddItem(shield)
	inv.AddItem(potion)
	inv.AddItem(bow)
	return inv, nil
}

func loadInventory(filename string) (*inventory.Inventory, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("открытие файла инвентаря: %w", err)
	}
	defer file.Close()
	inv := &inventory.Inventory{}
	if err := inv.Load(file); err != nil {
		return nil, fmt.Errorf("загрузка инвентаря: %w", err)
	}
	return inv, nil
}
