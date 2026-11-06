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

type inventoryStore interface {
	Load(filename string) ([]items.Item, error)
	Save(filename string, itemList []items.Item) error
}

type commandRunner struct {
	store inventoryStore
}

type commandOptions struct {
	command  string
	filename string
	name     string
	weight   float64
	id       int
}

func parseCommand(args []string, flagOutput io.Writer) (commandOptions, error) {
	var options commandOptions
	if len(args) == 0 {
		return options, errors.New("использование: inventory <list|add|use|delete> [флаги]")
	}

	options.command = args[0]
	flags := flag.NewFlagSet(options.command, flag.ContinueOnError)
	flags.SetOutput(flagOutput)
	flags.StringVar(&options.filename, "file", "inventory.txt", "файл инвентаря")
	switch options.command {
	case "list":
	case "add":
		flags.StringVar(&options.name, "name", "", "название предмета")
		flags.Float64Var(&options.weight, "weight", 0, "вес предмета")
	case "use", "delete":
		flags.IntVar(&options.id, "id", 0, "ID предмета")
	default:
		return options, fmt.Errorf("неизвестная команда %q", options.command)
	}
	if err := flags.Parse(args[1:]); err != nil {
		return options, err
	}
	if flags.NArg() != 0 || options.filename == "" {
		return options, errors.New("лишние аргументы или пустой путь к файлу")
	}
	switch options.command {
	case "add":
		if strings.ContainsAny(options.name, "|\r\n") {
			return options, errors.New("название предмета содержит разделитель или перенос строки")
		}
		options.name = strings.TrimSpace(options.name)
		if options.name == "" || options.weight <= 0 || math.IsNaN(options.weight) || math.IsInf(options.weight, 0) {
			return options, errors.New("для add нужны допустимое -name и положительное конечное -weight")
		}
	case "use", "delete":
		if options.id <= 0 {
			return options, errors.New("укажи положительный -id")
		}
	}
	return options, nil
}

func (runner commandRunner) run(args []string, stdout, stderr io.Writer) int {
	var flagOutput strings.Builder
	options, err := parseCommand(args, &flagOutput)
	if errors.Is(err, flag.ErrHelp) {
		if _, err := io.WriteString(stderr, flagOutput.String()); err != nil {
			return 1
		}
		return 0
	}
	if err != nil {
		if flagOutput.Len() > 0 {
			return exitWithDiagnostic(stderr, 2, flagOutput.String())
		}
		return exitWithDiagnostic(stderr, 2, err.Error()+"\n")
	}

	itemList, err := runner.store.Load(options.filename)
	if err != nil {
		return reportError(stderr, err)
	}
	var inv inventory.Inventory
	if err := inv.LoadItems(itemList); err != nil {
		return reportError(stderr, err)
	}
	var message string
	switch options.command {
	case "list":
		for _, item := range inv.Items() {
			if _, err := fmt.Fprintln(stdout, items.DescribeItem(item)); err != nil {
				return reportError(stderr, fmt.Errorf("вывод списка предметов: %w", err))
			}
		}
		return 0
	case "add":
		item := items.NewPotion(0, options.name, options.weight, "без эффекта", 1)
		itemID := inv.AddItem(item)
		if itemID == 0 {
			return reportError(stderr, errors.New("не удалось назначить ID предмету"))
		}
		message = fmt.Sprintf("Добавлен предмет %q с ID %d", item.Name(), itemID)
	case "use":
		item, found := inv.FindItemByID(options.id)
		if !found {
			return reportError(stderr, fmt.Errorf("предмет с ID %d не найден", options.id))
		}
		message, err = items.UseDangerousItem(item)
		if err != nil {
			return reportError(stderr, err)
		}
	case "delete":
		if !inv.DeleteItem(options.id) {
			return reportError(stderr, fmt.Errorf("предмет с ID %d не найден", options.id))
		}
		message = fmt.Sprintf("Удалён предмет с ID %d", options.id)
	}
	if err := runner.store.Save(options.filename, inv.Items()); err != nil {
		return reportError(stderr, err)
	}
	if _, err := fmt.Fprintln(stdout, message); err != nil {
		return reportError(stderr, fmt.Errorf("изменения сохранены, вывод результата: %w", err))
	}
	return 0
}

func reportError(stderr io.Writer, err error) int {
	var parseErr *storage.ParseItemError
	switch {
	case errors.As(err, &parseErr):
		return exitWithDiagnostic(stderr, 1, fmt.Sprintf("Повреждённая строка %d (%q): %v\n", parseErr.Line, parseErr.Value, parseErr.Err))
	case errors.Is(err, items.ErrBroken):
		return exitWithDiagnostic(stderr, 1, fmt.Sprintf("Предмет сломан: %v\n", err))
	case errors.Is(err, items.ErrExhausted):
		return exitWithDiagnostic(stderr, 1, fmt.Sprintf("Предмет исчерпан: %v\n", err))
	default:
		return exitWithDiagnostic(stderr, 1, err.Error()+"\n")
	}
}

func exitWithDiagnostic(stderr io.Writer, code int, message string) int {
	// Если stderr недоступен, код выхода всё равно сообщает исходную ошибку.
	_, _ = io.WriteString(stderr, message)
	return code
}
