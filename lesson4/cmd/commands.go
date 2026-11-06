package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"strings"

	"gostudy/lesson4/internal/inventory"
	"gostudy/lesson4/internal/items"
	"gostudy/lesson4/internal/storage"
)

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "Использование: inventory <list|add|use|delete> [флаги]")
		return 2
	}

	command := args[0]
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	filename := flags.String("file", "inventory.txt", "файл инвентаря")
	var name *string
	var weight *float64
	var id *int
	switch command {
	case "list":
	case "add":
		name = flags.String("name", "", "название предмета")
		weight = flags.Float64("weight", 0, "вес предмета")
	case "use", "delete":
		id = flags.Int("id", 0, "ID предмета")
	default:
		fmt.Fprintf(stderr, "Неизвестная команда %q\n", command)
		return 2
	}
	if err := flags.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || *filename == "" {
		fmt.Fprintln(stderr, "Лишние аргументы или пустой путь к файлу")
		return 2
	}
	if command == "add" && (strings.TrimSpace(*name) == "" || strings.ContainsAny(*name, "|\r\n") ||
		*weight <= 0 || math.IsNaN(*weight) || math.IsInf(*weight, 0)) {
		fmt.Fprintln(stderr, "Для add нужны допустимое -name и положительное конечное -weight")
		return 2
	}
	if id != nil && *id <= 0 {
		fmt.Fprintln(stderr, "Укажи положительный -id")
		return 2
	}

	itemList, err := storage.LoadFile(*filename)
	if err != nil {
		return reportError(stderr, err)
	}
	inv := inventory.Inventory{Items: itemList}
	var message string
	switch command {
	case "list":
		for _, item := range inv.Items {
			fmt.Fprintln(stdout, items.DescribeItem(item))
		}
		return 0
	case "add":
		item := items.NewPotion(0, strings.TrimSpace(*name), *weight, "без эффекта", 1)
		itemID := inv.AddItem(item)
		if itemID == 0 {
			return reportError(stderr, fmt.Errorf("не удалось назначить ID предмету"))
		}
		message = fmt.Sprintf("Добавлен предмет %q с ID %d", item.Name(), itemID)
	case "use":
		item, found := inv.FindItemByID(*id)
		if !found {
			return reportError(stderr, fmt.Errorf("предмет с ID %d не найден", *id))
		}
		message, err = items.UseDangerousItem(item)
		if err != nil {
			return reportError(stderr, err)
		}
	case "delete":
		if !inv.DeleteItem(*id) {
			return reportError(stderr, fmt.Errorf("предмет с ID %d не найден", *id))
		}
		message = fmt.Sprintf("Удалён предмет с ID %d", *id)
	}
	if err := storage.SaveFile(*filename, inv.Items); err != nil {
		return reportError(stderr, err)
	}
	fmt.Fprintln(stdout, message)
	return 0
}

func reportError(stderr io.Writer, err error) int {
	var parseErr *storage.ParseItemError
	switch {
	case errors.As(err, &parseErr):
		fmt.Fprintf(stderr, "Повреждённая строка %d (%q): %v\n", parseErr.Line, parseErr.Value, parseErr.Err)
	case errors.Is(err, items.ErrBroken):
		fmt.Fprintln(stderr, "Предмет сломан:", err)
	case errors.Is(err, items.ErrExhausted):
		fmt.Fprintln(stderr, "Предмет исчерпан:", err)
	default:
		fmt.Fprintln(stderr, err)
	}
	return 1
}
