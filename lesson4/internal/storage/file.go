package storage

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"gostudy/lesson4/internal/items"
)

func Save(w io.Writer, itemList []items.Item) error {
	for _, item := range itemList {
		if items.IsNilItem(item) {
			return fmt.Errorf("сохранение: %w", items.ErrNilItem)
		}
		storable, ok := item.(items.Storable)
		if !ok {
			return &items.ItemError{ID: items.IDOf(item), Name: item.Name(), Err: ErrNotStorable}
		}
		if err := storable.Serialize(w); err != nil {
			return &items.ItemError{ID: items.IDOf(item), Name: item.Name(), Err: err}
		}
	}
	return nil
}

// Load останавливается на первой повреждённой строке и не возвращает частичный инвентарь.
func Load(r io.Reader) ([]items.Item, error) {
	scanner := bufio.NewScanner(r)
	loaded := make([]items.Item, 0)
	line := 0
	for scanner.Scan() {
		line++
		record := scanner.Text()
		item, err := items.ReadItem(strings.NewReader(record))
		if err != nil {
			return nil, &ParseItemError{Line: line, Value: record, Err: err}
		}
		loaded = append(loaded, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("чтение строки %d: %w", line+1, err)
	}
	return loaded, nil
}

func LoadFile(filename string) ([]items.Item, error) {
	file, err := os.Open(filename)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("открытие инвентаря %q: %w", filename, err)
	}
	itemList, loadErr := Load(file)
	closeErr := file.Close()
	if loadErr != nil {
		loadErr = fmt.Errorf("загрузка инвентаря %q: %w", filename, loadErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("закрытие инвентаря %q: %w", filename, closeErr)
	}
	if err := errors.Join(loadErr, closeErr); err != nil {
		return nil, err
	}
	return itemList, nil
}

func SaveFile(filename string, itemList []items.Item) error {
	var data bytes.Buffer
	if err := Save(&data, itemList); err != nil {
		return fmt.Errorf("подготовка инвентаря %q: %w", filename, err)
	}
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("открытие инвентаря для записи %q: %w", filename, err)
	}
	_, writeErr := data.WriteTo(file)
	closeErr := file.Close()
	if writeErr != nil {
		writeErr = fmt.Errorf("запись инвентаря %q: %w", filename, writeErr)
	}
	if closeErr != nil {
		closeErr = fmt.Errorf("закрытие инвентаря %q: %w", filename, closeErr)
	}
	return errors.Join(writeErr, closeErr)
}
